package context

import (
	"mmo/common/cmd"
	"mmo/servers/scene/define"
)

type CmdContext struct {
	Player define.IPlayer
}

func NewCmdContext(player define.IPlayer) cmd.IContext {
	return &CmdContext{
		Player: player,
	}
}
