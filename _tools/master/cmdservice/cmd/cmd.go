package cmd

import (
	"strings"
)

type DoFunc func(ctx IContext, args []string, cb func(string))

type IContext interface{}

type Mgr struct {
	cmds map[string]DoFunc
}

func NewMgr() *Mgr {
	return &Mgr{
		cmds: make(map[string]DoFunc),
	}
}

func (m *Mgr) Register(op string, doFunc DoFunc) {
	m.cmds[op] = doFunc
}

func (m *Mgr) Call(ctx IContext, cmd string, cb func(result string)) {
	if cmd == "" {
		return
	}
	subs := strings.Split(cmd, " ")

	f := m.cmds[subs[0]]
	if f == nil {
		cb("no this cmd")
		return
	}

	f(ctx, subs[1:], cb)
}
