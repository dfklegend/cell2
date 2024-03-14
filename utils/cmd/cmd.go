package cmd

/*
	cmd包
	提供一个简单的处理输入字符串的控制台循环处理机制
	工程中可以使用来处理一些简单的控制台输入
*/

import (
	"strings"

	"github.com/dfklegend/cell2/utils/logger"
)

type ICmd interface {
	GetName() string
	Do(args []string)
}

type Mgr struct {
	items map[string]ICmd
}

func NewCmdMgr() *Mgr {
	return &Mgr{items: make(map[string]ICmd)}
}

func (s *Mgr) Register(c ICmd) {
	logger.Log.Infof("register cmd:%s", c.GetName())
	s.items[(c).GetName()] = c
}

func (s *Mgr) Dispatch(content string) {
	subs := strings.Split(content, " ")
	if len(subs) < 1 {
		return
	}
	name := subs[0]
	c, ok := s.items[name]
	if !ok {
		return
	}
	c.Do(subs[1:])
}
