package etcd

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dfklegend/cell2/node/cluster"
)

type TestCluster struct {
	Address  string
	Name     string
	ID       string
	Services []string
	CB       TestCBFunc
}

type TestCBFunc func(*TestCluster, []*cluster.Member)

func (t *TestCluster) GetAddress() string {
	return t.Address
}

func (t *TestCluster) GetName() string {
	return t.Name
}

func (t *TestCluster) GetID() string {
	return t.ID
}

func (t *TestCluster) GetState() int {
	return 0
}

func (t *TestCluster) GetServices() []string {
	return t.Services
}

func (t *TestCluster) UpdateClusterTopology(members []*cluster.Member) {
	log.Printf("%v UpdateClusterTopology %v\n",
		t.ID, len(members))
	for i := 0; i < len(members); i++ {
		one := members[i]
		log.Println(one)
	}
	if t.CB != nil {
		t.CB(t, members)
	}
}

// 	需要开启etcd
//	测试加入节点
func TestStartMember(t *testing.T) {
	if testing.Short() {
		return
	}

	c := &TestCluster{
		Address: "127.0.0.1:1000",
		Name:    "unittest",
		ID:      "001",
		Services: []string{
			"chat.chat1",
			"chat.chat2",
		},
	}

	p := newProvider(c)

	var finalMembers []*cluster.Member
	c.CB = func(c *TestCluster, members []*cluster.Member) {
		finalMembers = members
	}
	defer p.Shutdown(true)
	time.Sleep(2 * time.Second)
	assert.Equal(t, 1, len(finalMembers))
}

func newProvider(c cluster.ICluster) *Provider {
	p, _ := newTestProvider()
	p.StartMember(c)
	return p
}

// 	测试节点移除
func TestRemoveMember(t *testing.T) {
	if testing.Short() {
		return
	}

	cluster1 := &TestCluster{
		Address: "127.0.0.1:1000",
		Name:    "unittest",
		ID:      "001",
		Services: []string{
			"chat.chat1",
			"chat.chat2",
		},
	}
	var f1 []*cluster.Member
	cluster1.CB = func(c *TestCluster, members []*cluster.Member) {
		f1 = members
	}

	p1 := newProvider(cluster1)

	cluster2 := &TestCluster{
		Address: "127.0.0.1:1001",
		Name:    "unittest",
		ID:      "002",
		Services: []string{
			"gate.gate1",
			"gate.gate2",
		},
	}
	var f2 []*cluster.Member
	cluster2.CB = func(c *TestCluster, members []*cluster.Member) {
		f2 = members
	}

	p2 := newProvider(cluster2)

	time.Sleep(2 * time.Second)
	assert.Equal(t, 2, len(f1))
	assert.Equal(t, 2, len(f2))

	p2.Shutdown(true)
	plog.Info("stop p2")
	time.Sleep(1 * time.Second)
	assert.Equal(t, 1, len(f1))

	p1.Shutdown(true)
	plog.Info("stop p1")
}

// 	测试中断流程后，能否正常恢复
func TestBadKeepAliveProcess(t *testing.T) {
	if testing.Short() {
		return
	}

	cluster1 := &TestCluster{
		Address: "127.0.0.1:1000",
		Name:    "unittest",
		ID:      "001",
		Services: []string{
			"chat.chat1",
			"chat.chat2",
		},
	}
	var f1 []*cluster.Member
	cluster1.CB = func(c *TestCluster, members []*cluster.Member) {
		f1 = members
	}

	p1 := newProvider(cluster1)

	cluster2 := &TestCluster{
		Address: "127.0.0.1:1001",
		Name:    "unittest",
		ID:      "002",
		Services: []string{
			"gate.gate1",
			"gate.gate2",
		},
	}
	var f2 []*cluster.Member
	cluster2.CB = func(c *TestCluster, members []*cluster.Member) {
		f2 = members
	}

	p2 := newProvider(cluster2)
	p1.debugShowEvent = true
	p2.debugShowEvent = true

	time.Sleep(1 * time.Second)
	log.Println("set debugBadLease true")
	p1.debugBadLease = true

	time.Sleep(5 * time.Second)
	// 自己不会删除，所以还是2
	assert.Equal(t, 2, len(f1))
	// p1没了
	assert.Equal(t, 1, len(f2))

	log.Println("set debugBadLease false")
	p1.debugBadLease = false
	time.Sleep(5 * time.Second)
	assert.Equal(t, 2, len(f1))
	assert.Equal(t, 2, len(f2))
	p1.Shutdown(true)
	p2.Shutdown(true)
}

func TestUpdateState(t *testing.T) {
	if testing.Short() {
		return
	}

	c := &TestCluster{
		Address: "127.0.0.1:1000",
		Name:    "unittest",
		ID:      "001",
		Services: []string{
			"chat.chat1",
			"chat.chat2",
		},
	}

	p := newProvider(c)
	p.debugShowEvent = true

	var finalMembers []*cluster.Member
	c.CB = func(c *TestCluster, members []*cluster.Member) {
		finalMembers = members
	}
	defer p.Shutdown(true)

	time.Sleep(2 * time.Second)
	log.Printf("state -> 1")
	p.UpdateClusterState(1)
	time.Sleep(2 * time.Second)
	log.Printf("state -> 2")
	p.UpdateClusterState(2)
	time.Sleep(2 * time.Second)

	assert.Equal(t, 1, len(finalMembers))
	assert.Equal(t, 2, finalMembers[0].State)
}

// 测试获取组，不会获取到通配的组
// 比如 card, card1
// 获取card会得到card,card1的信息
func TestClusterNameProcess(t *testing.T) {
	if testing.Short() {
		return
	}

	// 组1
	cluster1 := &TestCluster{
		Address: "127.0.0.1:1000",
		Name:    "unittest",
		ID:      "001",
		Services: []string{
			"chat.chat1",
			"chat.chat2",
		},
	}
	var f1 []*cluster.Member
	cluster1.CB = func(c *TestCluster, members []*cluster.Member) {
		f1 = members
	}

	p1 := newProvider(cluster1)

	cluster2 := &TestCluster{
		Address: "127.0.0.1:1001",
		Name:    "unittest",
		ID:      "002",
		Services: []string{
			"gate.gate1",
			"gate.gate2",
		},
	}
	var f2 []*cluster.Member
	cluster2.CB = func(c *TestCluster, members []*cluster.Member) {
		f2 = members
	}

	p2 := newProvider(cluster2)

	// 组2
	cluster3 := &TestCluster{
		Address: "127.0.0.1:1001",
		Name:    "unittest3",
		ID:      "003",
		Services: []string{
			"gate.gate1",
			"gate.gate2",
		},
	}
	var f3 []*cluster.Member
	cluster3.CB = func(c *TestCluster, members []*cluster.Member) {
		f3 = members
	}

	p3 := newProvider(cluster3)

	p1.debugShowEvent = true
	p2.debugShowEvent = true
	p3.debugShowEvent = true

	time.Sleep(1 * time.Second)
	log.Printf("f1: %v\n", len(f1))
	log.Printf("f2: %v\n", len(f2))
	log.Printf("f3: %v\n", len(f3))

	time.Sleep(5 * time.Second)
	log.Printf("f1: %v\n", len(f1))
	log.Printf("f2: %v\n", len(f2))
	log.Printf("f3: %v\n", len(f3))

	assert.Equal(t, 2, len(f1))
	assert.Equal(t, 1, len(f3))
}
