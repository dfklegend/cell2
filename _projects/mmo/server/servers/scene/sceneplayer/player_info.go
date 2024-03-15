package sceneplayer

import (
	mymsg "mmo/messages"
)

// LoadInfo
// 载入玩家数据
func (p *ScenePlayer) LoadInfo(info *mymsg.PlayerInfo) {
	// 每个system拉取数据
	p.systemsLoadData(info)
}

// InitNewPlayerInfo 初始化新建角色
func (p *ScenePlayer) InitNewPlayerInfo() {
	p.systemsInitData()
	p.SetDirt()
}

func (p *ScenePlayer) MakeSaveInfo() *mymsg.PlayerInfo {
	info := &mymsg.PlayerInfo{}
	p.systemsSaveData(info)
	return info
}

func (p *ScenePlayer) ReqCharInfo() {
	p.systemsReqInfo()
}
