package test1

import (
	"reflect"

	"mmo/common/define"
	"mmo/common/factory"
	mymsg "mmo/messages"
	"mmo/messages/cproto"
	"mmo/modules/systems"
	sysfactory "mmo/servers/scene/sceneplayer/systems/factory"
)

func init() {
	name := "test1"
	sysfactory.GetFactory().RegisterFunc(name, func(args ...any) factory.IObject {
		return newExample(name)
	})
}

func Visit() {}

type Test1 struct {
	*systems.BaseSystem
}

func newExample(name string) systems.ISystem {
	return &Test1{
		BaseSystem: systems.NewBaseSystem(name),
	}
}

func (t *Test1) OnCreate() {
	t.bindEvents(true)

	t.RegisterCmdHandler("add", reflect.TypeOf(&cproto.TestAdd{}), t.onHandleAdd)

	logger := t.GetPlayer().GetNodeService().GetLogger()
	logger.Infof("Test1.OnCreate")
}

func (t *Test1) OnDestroy() {
}

func (t *Test1) bindEvents(bind bool) {
	//events := t.player.GetEvents()
}

func (t *Test1) PushInfoToClient() {
	t.PushCmd("test1", &cproto.TestAddRet{
		Result: 101,
	})
	t.PushCmd("test2", &cproto.TestAddRet{
		Result: 102,
	})
}

func (t *Test1) LoadData(info *mymsg.PlayerInfo) {
}

func (t *Test1) SaveData(info *mymsg.PlayerInfo) {
}

func (t *Test1) onHandleAdd(srcArgs any, cb func(ret any, code int32)) {
	args, _ := srcArgs.(*cproto.TestAdd)
	if args == nil {
		return
	}
	cb(&cproto.TestAddRet{
		Result: args.I + args.J,
	}, int32(define.Succ))
}

func (t *Test1) Hello() {
	logger := t.GetPlayer().GetNodeService().GetLogger()
	logger.Infof("hello from test1")
}
