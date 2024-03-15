package fightmode

import (
	"mmo/messages/cproto"
	"mmo/servers/logic_old/define"
)

// ISceneFightMode
// 战斗模式, 不同场景可能是不同的战斗模式
// 这些战斗模式需要准备自己的战斗数据，发送给场景
// 和玩家有关

type ISceneFightData interface {
}

type ISceneFightResult interface {
}

type ISceneFightMode interface {

	// InitData 设置数据
	InitData(player define.ILogicPlayer, data ISceneFightData)
	SendInitDataToClient(cb func())

	// InitSceneInfo after scene alloced
	InitSceneInfo(sceneService string, sceneId uint64, token int32)
	SendInitDataToScene(cb func())

	OnFightResult(result ISceneFightResult)
}

type CardFightData struct {
	Cards []*cproto.CharCard
}
