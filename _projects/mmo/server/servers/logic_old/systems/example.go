package systems

//
//import (
//	"reflect"
//
//	"mmo/common/define"
//	"mmo/common/factory"
//	mymsg "mmo/messages"
//	"mmo/messages/cproto"
//	define2 "mmo/servers/logic_old/define"
//	"mmo/servers/logic_old/fightmode"
//	sysfactory "mmo/servers/logic_old/systems/factory"
//)
//
//func init() {
//	sysfactory.GetFactory().RegisterFunc("example", func(args ...any) factory.IObject {
//		return newExample("example")
//	})
//}
//
//type Example struct {
//	*BaseSystem
//
//	curFight fightmode.ISceneFightMode
//}
//
//func newExample(name string) define2.ISystem {
//	return &Example{
//		BaseSystem: NewBaseSystem(name),
//	}
//}
//
//func (c *Example) OnCreate() {
//	c.bindEvents(true)
//
//	c.RegisterCmdHandler("add", reflect.TypeOf(&cproto.TestAdd{}), c.onHandleAdd)
//}
//
//func (c *Example) OnDestroy() {
//}
//
//func (c *Example) bindEvents(bind bool) {
//	//events := c.player.GetEvents()
//}
//
//func (c *Example) PushInfoToClient() {
//	c.PushCmd("test1", &cproto.TestAddRet{
//		Result: 101,
//	})
//	c.PushCmd("test2", &cproto.TestAddRet{
//		Result: 102,
//	})
//}
//
//func (c *Example) LoadData(info *mymsg.PlayerInfo) {
//}
//
//func (c *Example) SaveData(info *mymsg.PlayerInfo) {
//}
//
//func (c *Example) onHandleAdd(srcArgs any, cb func(ret any, code int32)) {
//	args, _ := srcArgs.(*cproto.TestAdd)
//	if args == nil {
//		return
//	}
//	cb(&cproto.TestAddRet{
//		Result: args.I + args.J,
//	}, int32(define.Succ))
//}
