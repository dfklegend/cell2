package scenecfg

import (
	"fmt"
)

var (
	cfgs *Mgr
)

type Mgr struct {
	path string
	cfgs map[int]*SceneCfg
}

func NewMgr(path string) *Mgr {
	return &Mgr{
		cfgs: make(map[int]*SceneCfg),
		path: path,
	}
}

func InitMgr(path string) *Mgr {
	cfgs = NewMgr(path)
	return cfgs
}

func GetMgr() *Mgr {
	return cfgs
}

func (m *Mgr) Load(cfgId int) {
	_, ok := m.cfgs[cfgId]
	if ok {
		return
	}
	cfg := LoadSceneCfg(m.path, fmt.Sprintf("%v", cfgId))
	if cfg.Monsters == nil && cfg.Exits == nil {
		return
	}
	m.cfgs[cfgId] = cfg
}

func (m *Mgr) GetItem(cfgId int) *SceneCfg {
	return m.cfgs[cfgId]
}
