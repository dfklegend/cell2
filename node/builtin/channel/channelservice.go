package channel

import (
	"fmt"
	"sync"

	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/common"
)

var (
	tempIdService = common.NewSerialIdService()
)

// Service 频道服务
// 可以向固定的成员推送消息
// 可以被具体Service集成
// 目前并不需要锁(每个service一个独立的channelservice)
type Service struct {
	channels sync.Map
	ns       *service.NodeService
}

func NewChannelService(ns *service.NodeService) *Service {
	return &Service{
		ns: ns,
	}
}

func (s *Service) AddChannel(name string) *Channel {
	c, ok := s.channels.Load(name)
	if ok {
		return c.(*Channel)
	}
	nc := NewChannel(name, s.ns)
	s.channels.Store(name, nc)
	return nc
}

func (s *Service) GetChannel(name string) *Channel {
	c, ok := s.channels.Load(name)
	if !ok {
		return nil
	}
	return c.(*Channel)
}

func (s *Service) DeleteChannel(name string) {
	s.channels.Delete(name)
}

func (s *Service) AddToChannel(name string, serverId string,
	netId uint32) *Channel {
	c := s.AddChannel(name)
	c.Add(serverId, netId)
	return c
}

func (s *Service) LeaveFromChannel(name string, serverId string,
	netId uint32) {
	c := s.GetChannel(name)
	if c != nil {
		c.Leave(serverId, netId)
	}
}

// PushMessageById 直接向[frontId,netId]推送消息
func (s *Service) PushMessageById(ns *service.NodeService, serverId string, netId uint32,
	route string, msg interface{}, cbFunc apientry.HandlerCBFunc) {
	PushMessageById(ns, serverId, uint32(netId), route, msg)
	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
}

func (s *Service) PushMessageByIds(ns *service.NodeService, serverId string, netIds []uint32,
	route string, msg interface{}, cbFunc apientry.HandlerCBFunc) {
	PushMessageByIds(ns, serverId, netIds, route, msg)
	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
}

// AllocTempChannel 可以申请一个临时的channel，用来方便发送请求
func (s *Service) AllocTempChannel() *Channel {
	name := fmt.Sprintf("_temp_chanel_%v", tempIdService.AllocId())
	return s.AddChannel(name)
}

// c.Add..
// c.PushMessage

func (s *Service) FreeTempChannel(c *Channel) {
	s.DeleteChannel(c.GetName())
}
