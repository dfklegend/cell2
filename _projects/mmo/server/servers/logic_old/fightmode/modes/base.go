package modes

import (
	"mmo/servers/logic_old/define"
	"mmo/servers/logic_old/fightmode"
)

type BaseMode struct {
	player define.ILogicPlayer
	uid    int64

	sceneService string
	sceneId      uint64
	token        int32
}

func newBaseMode() *BaseMode {
	return &BaseMode{}
}

func (m *BaseMode) InitData(player define.ILogicPlayer, data fightmode.ISceneFightData) {
	m.player = player
	m.uid = player.GetUId()
}

func (m *BaseMode) SendInitDataToClient(cb func()) {
	cb()
}

// InitSceneInfo 设置分配到的场景参数
func (m *BaseMode) InitSceneInfo(sceneService string, sceneId uint64, token int32) {
	m.sceneService = sceneService
	m.sceneId = sceneId
	m.token = token
}

func (m *BaseMode) SendInitDataToScene(cb func()) {
	cb()
}
