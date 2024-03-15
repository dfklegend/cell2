package test2

import (
	"reflect"

	"github.com/dfklegend/cell2/utils/event/light"

	"mmo/common/define"
	"mmo/common/factory"
	mymsg "mmo/messages"
	"mmo/messages/cproto"
	"mmo/modules/systems"
	sysfactory "mmo/servers/scene/sceneplayer/systems/factory"
	"mmo/servers/scene/sceneplayer/systems/interfaces"
)

func init() {
	name := "test2"
	sysfactory.GetFactory().RegisterFunc(name, func(args ...any) factory.IObject {
		return newExample(name)
	})
}

func Visit() {}

type Test2 struct {
	*systems.BaseSystem
}

func newExample(name string) systems.ISystem {
	return &Test2{
		BaseSystem: systems.NewBaseSystem(name),
	}
}

func (t *Test2) OnCreate() {
	t.bindEvents(true)

	t.RegisterCmdHandler("add", reflect.TypeOf(&cproto.TestAdd{}), t.onHandleAdd)

	logger := t.GetPlayer().GetNodeService().GetLogger()
	logger.Infof("Test2.OnCreate")

}

func (t *Test2) OnDestroy() {
	t.bindEvents(false)
}

func (t *Test2) bindEvents(bind bool) {
	events := t.GetPlayer().GetEvents()
	light.BindEventWithReceiver(bind, events, "onOffline", t, t.onOffline)
}

func (t *Test2) PushInfoToClient() {
	t.PushCmd("test1", &cproto.TestAddRet{
		Result: 101,
	})
	t.PushCmd("test2", &cproto.TestAddRet{
		Result: 102,
	})
}

func (t *Test2) LoadData(info *mymsg.PlayerInfo) {
}

func (t *Test2) SaveData(info *mymsg.PlayerInfo) {
}

func (t *Test2) onHandleAdd(srcArgs any, cb func(ret any, code int32)) {
	args, _ := srcArgs.(*cproto.TestAdd)
	if args == nil {
		return
	}
	cb(&cproto.TestAddRet{
		Result: args.I + args.J,
	}, int32(define.Succ))
}

func (t *Test2) Hello() {
	logger := t.GetPlayer().GetNodeService().GetLogger()
	t1 := t.GetPlayer().GetSystem("test1").(interfaces.Test1)

	logger.Infof("call test1.Hello")
	t1.Hello()
}

func (t *Test2) onOffline(args ...any) {
	logger := t.GetPlayer().GetNodeService().GetLogger()
	logger.Infof("test2 onOffline")
}
