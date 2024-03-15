package scenelogics

import (
	"mmo/servers/scene/define"
	logic2 "mmo/servers/scene/logic"
)

func init() {
	logic2.GetLogicFactory().Register("empty", func() logic2.ISceneLogic {
		//return newEmptyLogic()
		return nil
	})
}

type EmptyLogic struct {
	logic2.ISceneLogic
}

func newEmptyLogic() *EmptyLogic {
	return &EmptyLogic{}
}

func (l *EmptyLogic) Init(define.IScene) {
}

func (l *EmptyLogic) Start() {
}

func (l *EmptyLogic) IsOver() bool {
	return false
}

func (l *EmptyLogic) Update() {
}

func (l *EmptyLogic) Destroy() {
}
