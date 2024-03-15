package control

import (
	"reflect"

	"github.com/dfklegend/cell2/utils/logger/interfaces"

	"mmo/common/define"
	"mmo/common/factory"
	"mmo/messages/cproto"
	"mmo/modules/systems"
	scenedefine "mmo/servers/scene/define"
	entitydefine "mmo/servers/scene/entity/define"
	sysfactory "mmo/servers/scene/sceneplayer/systems/factory"
)

func init() {
	name := "control"
	sysfactory.GetFactory().RegisterFunc(name, func(args ...any) factory.IObject {
		return newSystem(name)
	})
}

func Visit() {}

// 处理玩家移动请求
type system struct {
	*systems.BaseSystem
	logger interfaces.Logger
}

func newSystem(name string) systems.ISystem {
	return &system{
		BaseSystem: systems.NewBaseSystem(name),
	}
}

func (s *system) OnCreate() {
	s.logger = s.GetPlayer().GetNodeService().GetLogger()
	s.logger.Infof("control.OnCreate")

	s.bindEvents(true)

	// cmds
	s.RegisterCmdHandler("moveto", reflect.TypeOf(&cproto.ReqMoveTo{}), s.onHandleMove)
	s.RegisterCmdHandler("stopmove", reflect.TypeOf(&cproto.ReqStopMove{}), s.onHandleStopMove)
}

func (s *system) bindEvents(bind bool) {
	//events := s.player.GetEvents()
}

func (s *system) onHandleMove(srcArgs any, cb func(ret any, code int32)) {
	//s.logger.Infof("control.onHandleMove")
	if srcArgs == nil {
		cb(nil, int32(define.ErrFaild))
		return
	}
	p := s.GetPlayer().(scenedefine.IPlayer)
	e := p.GetAvatarEntity()
	if e == nil {
		cb(nil, int32(define.ErrFaild))
		return
	}
	msg, ok := srcArgs.(*cproto.ReqMoveTo)
	if !ok {
		cb(nil, int32(define.ErrFaild))
		return
	}
	ctrl := e.GetComponent(entitydefine.Control).(entitydefine.IControl)
	ctrl.ReqMoveTo(msg.X, msg.Z)
	cb(nil, int32(define.Succ))
}

func (s *system) onHandleStopMove(srcArgs any, cb func(ret any, code int32)) {
	//s.logger.Infof("control.onHandleStopMove")
	if srcArgs == nil {
		cb(nil, int32(define.ErrFaild))
		return
	}
	p := s.GetPlayer().(scenedefine.IPlayer)
	e := p.GetAvatarEntity()
	if e == nil {
		cb(nil, int32(define.ErrFaild))
		return
	}
	msg, ok := srcArgs.(*cproto.ReqStopMove)
	if !ok {
		cb(nil, int32(define.ErrFaild))
		return
	}
	ctrl := e.GetComponent(entitydefine.Control).(entitydefine.IControl)
	ctrl.ReqStopMove(msg.X, msg.Z)
	cb(nil, int32(define.Succ))
}
