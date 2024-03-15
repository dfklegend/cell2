package systems

import (
	"reflect"

	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/serialize"
	"github.com/dfklegend/cell2/utils/serialize/proto"

	"mmo/common/define"
)

type CmdHandler func(args any, cb func(ret any, code int32))

type CmdItem struct {
	ArgsType reflect.Type
	Handler  CmdHandler
}

// Cmds
// system可以注册一些cmd，通过 logicHandler.SystemCmd来调用 避免定义太多外层接口
// 可以协议局部化
type Cmds struct {
	cmds       map[string]*CmdItem
	serializer serialize.Serializer
}

func newCmds() *Cmds {
	return &Cmds{
		cmds:       map[string]*CmdItem{},
		serializer: proto.GetDefaultSerializer(),
	}
}

func (c *Cmds) Register(cmd string, argsType reflect.Type, handler CmdHandler) {
	item := &CmdItem{
		ArgsType: argsType,
		Handler:  handler,
	}

	c.cmds[cmd] = item
}

func (c *Cmds) HasCmd(cmd string) bool {
	return c.cmds[cmd] != nil
}

func (c *Cmds) MakeArgs(item *CmdItem, bytes []byte) any {
	if item == nil {
		return nil
	}
	argType := item.ArgsType
	if argType == nil {
		return nil
	}
	arg := reflect.New(argType.Elem()).Interface()
	err := c.serializer.Unmarshal(bytes, arg)
	if err != nil {
		l.Log.Infof("arg Unmarshal [%v] failed: %v\n", string(bytes), err)
		return nil
	}
	return arg
}

func (c *Cmds) Request(cmd string, bytes []byte, cb func(ret []byte, code int32)) {
	item := c.cmds[cmd]
	if item == nil {
		return
	}

	args := c.MakeArgs(item, bytes)
	item.Handler(args, func(ret any, code int32) {
		var data []byte
		if code == int32(define.Succ) && ret != nil {
			d, err := c.serializer.Marshal(ret)
			if err != nil {
				l.L.Errorf("got error: %v", err)
			} else {
				data = d
			}
		}
		cb(data, code)
	})
}
