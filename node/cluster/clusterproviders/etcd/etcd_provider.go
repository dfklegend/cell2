package etcd

/*
	此provider代码从proto.actor的实现修改而来
	modify from proto.actor cluster/clusterproviders/etcd

	bug:
	如果keepalive失败，后续会由于lease失效，cluster更新异常
	发现是失败后，收到一个Delete事件，删除了自身的member
	添加时又由于是自身节点，忽略了
	修复: 删除时也添加判断自身即可
*/

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/dfklegend/cell2/node/cluster"
)

// 存储方式，节点id: json(Node)
type Provider struct {
	leaseID      clientv3.LeaseID
	cluster      cluster.ICluster
	baseKey      string
	clusterName  string
	deregistered bool
	shutdown     bool

	self *Node
	// 自身状态dirt
	selfStateDirt bool
	members       map[string]*Node // all, contains self.
	clusterError  error
	client        *clientv3.Client
	cancelWatch   func()
	cancelWatchCh chan bool
	keepAliveTTL  time.Duration
	retryInterval time.Duration
	revision      uint64
	// deregisterCritical time.Duration

	debugBadLease  bool
	debugShowEvent bool
}

func newTestProvider() (*Provider, error) {
	return NewWithConfig("/cell2", clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
}

func NewWithConfig(baseKey string, cfg clientv3.Config) (*Provider, error) {
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	p := &Provider{
		client:        client,
		keepAliveTTL:  3 * time.Second,
		retryInterval: 1 * time.Second,
		baseKey:       baseKey,
		members:       map[string]*Node{},
		cancelWatchCh: make(chan bool),
		selfStateDirt: false,
	}
	return p, nil
}

func (p *Provider) init(c cluster.ICluster) error {
	p.cluster = c
	addr := p.cluster.GetAddress()
	host, port, err := splitHostPort(addr)
	if err != nil {
		return err
	}

	p.cluster = c
	p.clusterName = p.cluster.GetName()
	// actorSystem.ID
	memberID := p.cluster.GetID()
	services := c.GetServices()

	// make self node infos
	nodeName := fmt.Sprintf("%v@%v", p.clusterName, memberID)
	p.self = NewNode(nodeName, host, port, services)
	p.self.SetState(c.GetState())
	p.self.SetMeta("id", p.getID())

	return nil
}

func (p *Provider) StartMember(c cluster.ICluster) error {
	if err := p.init(c); err != nil {
		return err
	}

	// fetch memberlist
	nodes, err := p.fetchNodes()
	if err != nil {
		return err
	}
	// initialize members
	p.updateNodesWithSelf(nodes)
	p.publishClusterTopologyEvent()
	p.startWatching()

	// register self
	if err := p.registerService(); err != nil {
		return err
	}
	ctx := context.TODO()
	p.startKeepAlive(ctx)
	return nil
}

func (p *Provider) StartClient(c cluster.ICluster) error {
	if err := p.init(c); err != nil {
		return err
	}
	nodes, err := p.fetchNodes()
	if err != nil {
		return err
	}
	// initialize members
	p.updateNodes(nodes)
	p.publishClusterTopologyEvent()
	p.startWatching()
	return nil
}

func (p *Provider) Shutdown(graceful bool) error {
	p.shutdown = true
	if !p.deregistered {
		err := p.deregisterService()
		if err != nil {
			plog.Error("deregisterMember", err)
			return err
		}
		p.deregistered = true
	}
	if p.cancelWatch != nil {
		p.cancelWatch()
		p.cancelWatch = nil
	}
	return nil
}

func (p *Provider) UpdateClusterState(state int) error {
	p.self.State = state
	p.selfStateDirt = true
	return nil
}

func (p *Provider) keepAliveForever(ctx context.Context) error {
	if p.self == nil {
		return fmt.Errorf("keepalive must be after initialize")
	}

	data, err := p.self.Serialize()
	if err != nil {
		return err
	}
	fullKey := p.getEtcdKey()
	if err != nil {
		return err
	}

	leaseId := p.getLeaseID()
	if leaseId <= 0 {
		newId, err := p.newLeaseID()
		if err != nil {
			return err
		}
		leaseId = newId
	}

	if leaseId <= 0 {
		return fmt.Errorf("grant lease failed. leaseId=%d", leaseId)
	}

	// 单元测试
	if p.debugBadLease {
		leaseId = 1234
	}

	_, err = p.client.Put(context.TODO(), fullKey, string(data), clientv3.WithLease(leaseId))
	if err != nil {
		return err
	}

	kaRespCh, err := p.client.KeepAlive(context.TODO(), leaseId)
	if err != nil {
		return err
	}

	for resp := range kaRespCh {
		if p.debugBadLease {
			// 重置掉之前的keepAlive
			p.client.Revoke(context.TODO(), leaseId)
			break
		}
		if resp == nil {
			return fmt.Errorf("keep alive failed. resp=%s", resp.String())
		}
		//plog.Info("", log.String("info", fmt.Sprintf("keep alive %s ttl=%d", p.getID(), resp.TTL)))
		if p.shutdown {
			return nil
		}

		if p.selfStateDirt {
			// 跳出，刷新数据
			p.selfStateDirt = false
			p.client.Revoke(context.TODO(), leaseId)
			return nil
		}
	}
	return nil
}

func (p *Provider) startKeepAlive(ctx context.Context) {
	go func() {
		for !p.shutdown {
			if err := ctx.Err(); err != nil {
				plog.Info("Keepalive was stopped.", err)
				return
			}

			if err := p.keepAliveForever(ctx); err != nil {
				plog.Info("Failure refreshing service TTL. ReTrying...", p.retryInterval)
				time.Sleep(p.retryInterval)
			}
		}
	}()
}

func (p *Provider) getID() string {
	return p.self.ID
}

func (p *Provider) getEtcdKey() string {
	return p.buildKey(p.clusterName, p.getID())
}

func (p *Provider) registerService() error {
	data, err := p.self.Serialize()
	if err != nil {
		return err
	}
	fullKey := p.getEtcdKey()
	if err != nil {
		return err
	}
	leaseId := p.getLeaseID()
	if leaseId <= 0 {
		_leaseId, err := p.newLeaseID()
		if err != nil {
			return err
		}
		leaseId = _leaseId
	}
	_, err = p.client.Put(context.TODO(), fullKey, string(data), clientv3.WithLease(leaseId))
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) deregisterService() error {
	fullKey := p.getEtcdKey()
	_, err := p.client.Delete(context.TODO(), fullKey)
	return err
}

func (p *Provider) handleWatchResponse(resp clientv3.WatchResponse) map[string]*Node {
	changes := map[string]*Node{}
	for _, ev := range resp.Events {
		key := string(ev.Kv.Key)
		nodeId, err := getNodeID(key, "/")
		if err != nil {
			plog.Error("Invalid member.", key)
			continue
		}

		if p.debugShowEvent {
			plog.Info("watch event", p.getID(), ev)
		}
		switch ev.Type {
		case clientv3.EventTypePut:
			node, err := NewNodeFromBytes(ev.Kv.Value)
			if err != nil {
				plog.Error("Invalid member.", key)
				continue
			}
			if p.self.Equal(node) {
				//plog.Debug("Skip add self.", log.String("key", key))
				plog.Debug("Skip add self.", key)
				continue
			}
			if _, ok := p.members[nodeId]; ok {
				//plog.Debug("Update member.", log.String("key", key))
				plog.Debug("Update member.", key)
			} else {
				//log.Debug("New member.", log.String("key", key))
				plog.Debug("New member.", key)
			}
			changes[nodeId] = node
		case clientv3.EventTypeDelete:
			node, ok := p.members[nodeId]
			if !ok {
				continue
			}

			// 自己节点也不要删
			if p.self.ID == node.ID {
				//plog.Debug("Skip delete self.", log.String("key", key))
				plog.Debug("Skip delete self.", key)
				continue
			}

			//plog.Debug("Delete member.", log.String("key", key))
			plog.Debug("Delete member.", key)
			cloned := *node
			cloned.SetAlive(false)
			changes[nodeId] = &cloned
		default:
			//plog.Error("Invalid etcd event.type.", log.String("key", key),
			//	log.String("type", ev.Type.String()))
			plog.Error("Invalid etcd event.type.", key,
				ev.Type.String())
		}
	}
	p.revision = uint64(resp.Header.GetRevision())
	return changes
}

func (p *Provider) keepWatching(ctx context.Context) error {
	clusterKey := p.buildClusterKey(p.clusterName)
	stream := p.client.Watch(ctx, clusterKey, clientv3.WithPrefix())
	return p._keepWatching(stream)
}

func (p *Provider) _keepWatching(stream clientv3.WatchChan) error {
	for resp := range stream {
		if err := resp.Err(); err != nil {
			plog.Error("Failure watching service.")
			return err
		}
		if len(resp.Events) <= 0 {
			//plog.Error("Empty etcd.events.", log.Int("events", len(resp.Events)))
			plog.Error("Empty etcd.events.", len(resp.Events))
			continue
		}
		nodesChanges := p.handleWatchResponse(resp)
		p.updateNodesWithChanges(nodesChanges)
		p.publishClusterTopologyEvent()
	}
	return nil
}

func (p *Provider) startWatching() {
	ctx := context.TODO()
	ctx, cancel := context.WithCancel(ctx)
	p.cancelWatch = cancel
	go func() {
		for !p.shutdown {
			if err := p.keepWatching(ctx); err != nil {
				//plog.Error("Failed to keepWatching.", log.Error(err))
				plog.Error("Failed to keepWatching.", err)
				p.clusterError = err
			}
		}
	}()
}

// GetHealthStatus returns an error if the cluster health status has problems
func (p *Provider) GetHealthStatus() error {
	return p.clusterError
}

func newContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), timeout)
}

func (p *Provider) buildKey(names ...string) string {
	return strings.Join(append([]string{p.baseKey}, names...), "/")
}

// 目录名后面加个/
// 防止前缀一致比如 card, card1, 获取card可能得到card,card1
func (p *Provider) buildClusterKey(clusterName string) string {
	key := p.buildKey(clusterName)
	key += "/"
	return key
}

func (p *Provider) fetchNodes() ([]*Node, error) {
	key := p.buildClusterKey(p.clusterName)
	resp, err := p.client.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var nodes []*Node
	for _, v := range resp.Kvs {
		n := Node{}
		if err := n.Deserialize(v.Value); err != nil {
			return nil, err
		}
		nodes = append(nodes, &n)
	}
	p.revision = uint64(resp.Header.GetRevision())
	// plog.Debug("fetch nodes",
	// 	log.Uint64("raft term", resp.Header.GetRaftTerm()),
	// 	log.Int64("revision", resp.Header.GetRevision()))
	return nodes, nil
}

func (p *Provider) updateNodes(members []*Node) {
	for _, n := range members {
		p.members[n.ID] = n
	}
}

func (p *Provider) updateNodesWithSelf(members []*Node) {
	p.updateNodes(members)
	p.members[p.self.ID] = p.self
}

func (p *Provider) updateNodesWithChanges(changes map[string]*Node) {
	//plog.Info("updateNodesWithChanges.", log.Int("changes", len(changes)))
	plog.Info("updateNodesWithChanges.", len(changes))
	for memberId, member := range changes {
		p.members[memberId] = member
		if !member.IsAlive() {
			delete(p.members, memberId)
		}
	}
}

func (p *Provider) createClusterTopologyEvent() []*cluster.Member {
	res := make([]*cluster.Member, len(p.members))
	i := 0
	for _, m := range p.members {
		res[i] = m.MemberStatus()
		i++
	}
	return res
}

func (p *Provider) publishClusterTopologyEvent() {
	res := p.createClusterTopologyEvent()
	//plog.Info("Update cluster.", log.Int("members", len(res)))
	plog.Info("Update cluster.", len(res))
	// for _, m := range res {
	// 	plog.Info("\t", log.Object("member", m))
	// }
	p.cluster.UpdateClusterTopology(res)
	// p.cluster.ActorSystem.EventStream.Publish(res)
}

func (p *Provider) getLeaseID() clientv3.LeaseID {
	val := (int64)(p.leaseID)
	return (clientv3.LeaseID)(atomic.LoadInt64(&val))
}

func (p *Provider) setLeaseID(leaseID clientv3.LeaseID) {
	val := (int64)(p.leaseID)
	atomic.StoreInt64(&val, (int64)(leaseID))
}

func (p *Provider) newLeaseID() (clientv3.LeaseID, error) {
	lease := clientv3.NewLease(p.client)
	ttlSecs := int64(p.keepAliveTTL / time.Second)
	resp, err := lease.Grant(context.TODO(), ttlSecs)
	if err != nil {
		return 0, err
	}
	return resp.ID, nil
}

func splitHostPort(addr string) (host string, port int, err error) {
	if h, p, e := net.SplitHostPort(addr); e != nil {
		if addr != "nonhost" {
			err = e
		}
		host = "nonhost"
		port = -1
	} else {
		host = h
		port, err = strconv.Atoi(p)
	}
	return
}
