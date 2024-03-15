package cmd

import (
	"fmt"
	"strings"

	"github.com/dfklegend/cell2/utils/logger"
)

/*
 * 提供一个方便的扩展基于字符串命令的形式
 * cmdName arg0 arg1 arg2 ...
 */

// IContext
// 传递上下文，调用者自己定义
type IContext interface {
}

type ICmd interface {
	GetName() string
	Do(ctx IContext, args []string, cb func(result string))
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

func (s *Mgr) Dispatch(ctx IContext, content string, cb func(result string)) {
	subs := strings.Split(content, " ")
	if len(subs) < 1 {
		return
	}
	name := subs[0]
	c, ok := s.items[name]
	if !ok {
		cb(fmt.Sprintf("can not find cmd: %v", name))
		return
	}
	c.Do(ctx, subs[1:], cb)
}
