package channel

import (
	"sync"

	"github.com/dfklegend/cell2/node/service"
)

// FrontGroup [frontId, netId]代表一个连接
type FrontGroup struct {
	NetIds []uint32
	mutex  sync.Mutex
}

func NewFrontGroup() *FrontGroup {
	return &FrontGroup{
		NetIds: make([]uint32, 0, 128),
	}
}

func (g *FrontGroup) Add(netId uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.NetIds = append(g.NetIds, netId)
}

func (g *FrontGroup) FindIndex(netId uint32) int {
	for i := 0; i < len(g.NetIds); i++ {
		if netId == g.NetIds[i] {
			return i
		}
	}
	return -1
}

func (g *FrontGroup) Remove(netId uint32) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	i := g.FindIndex(netId)
	// can not find
	if i == -1 {
		return
	}

	if i == 0 {
		g.NetIds = g.NetIds[1:]
		return
	}

	if i == len(g.NetIds)-1 {
		g.NetIds = g.NetIds[:i]
		return
	}

	g.NetIds = append(g.NetIds[:i], g.NetIds[i+1:]...)
}

func (g *FrontGroup) Lock() {
	g.mutex.Lock()
}

func (g *FrontGroup) Unlock() {
	g.mutex.Unlock()
}

type Channel struct {
	groups sync.Map
	name   string
	ns     *service.NodeService
}

func NewChannel(name string, ns *service.NodeService) *Channel {
	return &Channel{
		name: name,
		ns:   ns,
	}
}

func (c *Channel) GetName() string {
	return c.name
}

func (c *Channel) getGroup(frontId string, createIfMiss bool) *FrontGroup {
	g, ok := c.groups.Load(frontId)
	if ok {
		return g.(*FrontGroup)
	}
	if !createIfMiss {
		return nil
	}
	ng := NewFrontGroup()
	c.groups.Store(frontId, ng)
	return ng
}

// Add
/**
 * 加入频道
 * @param frontId{string} 前端服务器id
 * @param uid{uint32} 玩家对应的netId
 */
func (c *Channel) Add(frontId string, netId uint32) {
	g := c.getGroup(frontId, true)
	g.Add(netId)
}

func (c *Channel) Leave(frontId string, netId uint32) {
	g := c.getGroup(frontId, false)
	if g != nil {
		g.Remove(netId)
	}
}

func (c *Channel) Range(f func(k, v interface{}) bool) {
	c.groups.Range(f)
}

func (c *Channel) PushMessage(route string, msg any) {
	c.Range(func(k, v interface{}) bool {
		n := k.(string)
		g := v.(*FrontGroup)

		g.Lock()
		GetPushImpl().PushMessageByIds(c.ns, n, g.NetIds, route, msg)
		g.Unlock()
		return true
	})
}
