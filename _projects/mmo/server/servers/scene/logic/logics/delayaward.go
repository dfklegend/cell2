package scenelogics

import (
	"mmo/servers/scene/define"
	logic2 "mmo/servers/scene/logic"
)

func init() {
	logic2.GetLogicFactory().Register("delayaward", func() logic2.ISceneLogic {
		//return newDelayAwardLogic()
		return nil
	})
}

// DelayAwardLogic
// 玩家进入后，等待10s,给与奖励，场景结束
type DelayAwardLogic struct {
	logic2.ISceneLogic

	scene   define.IScene
	timeEnd int64
}

func newDelayAwardLogic() *DelayAwardLogic {
	return &DelayAwardLogic{}
}

func (d *DelayAwardLogic) Init(scene define.IScene) {
	d.scene = scene
}

func (l *DelayAwardLogic) Start() {
}

func (l *DelayAwardLogic) IsOver() bool {
	return false
}

func (d *DelayAwardLogic) Update() {
}

func (d *DelayAwardLogic) Destroy() {
}
