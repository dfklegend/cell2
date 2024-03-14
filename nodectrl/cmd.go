package nodectrl

import (
	"fmt"
)

type ICmd interface {
	Handle(n *NodeCtrl, args string) string
}

type Cmds struct {
	cmds map[string]ICmd
}

func NewCmds() *Cmds {
	return &Cmds{
		cmds: make(map[string]ICmd),
	}
}

func (c *Cmds) Register(op string, cmd ICmd) {
	c.cmds[op] = cmd
}

func (c *Cmds) Handle(n *NodeCtrl, op string, args string) string {
	cmd := c.cmds[op]
	if cmd == nil {
		return fmt.Sprintf("%v not support", op)
	}
	return cmd.Handle(n, args)
}
