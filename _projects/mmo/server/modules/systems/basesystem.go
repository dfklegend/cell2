package systems

import (
	"reflect"

	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/serialize/proto"

	mymsg "mmo/messages"
	"mmo/messages/cproto"
)

type BaseSystem struct {
	ISystem
	name   string
	cmds   *Cmds
	player ILogicPlayer
}

func NewBaseSystem(name string) *BaseSystem {
	return &BaseSystem{
		name: name,
		cmds: newCmds(),
	}
}

func (c *BaseSystem) Init(player ILogicPlayer) {
	c.player = player
}

func (c *BaseSystem) GetPlayer() ILogicPlayer {
	return c.player
}

func (c *BaseSystem) OnCreate() {
}

func (c *BaseSystem) OnDestroy() {
}

func (c *BaseSystem) PushInfoToClient() {
}

func (c *BaseSystem) OnEnterWorld(switchLine bool) {
}

func (c *BaseSystem) InitData() {
}

func (c *BaseSystem) LoadData(info *mymsg.PlayerInfo) {
}

func (c *BaseSystem) SaveData(info *mymsg.PlayerInfo) {
}

func (c *BaseSystem) RegisterCmdHandler(cmd string, argsType reflect.Type, handler CmdHandler) {
	c.cmds.Register(cmd, argsType, handler)
}

func (c *BaseSystem) Request(cmd string, args []byte, cb func(ret []byte, errCode int32)) {
	c.cmds.Request(cmd, args, cb)
}

func (c *BaseSystem) HasCmd(cmd string) bool {
	return c.cmds.HasCmd(cmd)
}

// PushCmd
// 推送系统cmd
func (c *BaseSystem) PushCmd(cmd string, args any) {
	serializer := proto.GetDefaultSerializer()
	bytes, err := serializer.Marshal(args)
	if err != nil {
		l.L.Errorf("PushCmd Marshal error! err: %v", err)
		return
	}
	c.player.PushMsg("servercmd", &cproto.ServerSystemCmd{
		System: c.name,
		Cmd:    cmd,
		Args:   bytes,
	})
}
