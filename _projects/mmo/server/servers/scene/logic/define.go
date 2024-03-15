package logic

import (
	"mmo/common/entity"
	"mmo/servers/scene/define"
)

// ISceneLogic
// 场景逻辑驱动
type ISceneLogic interface {
	// 创建

	Init(define.IScene)
	CreateSpace() define.ISpace
	SetSceneInitData(data any)

	// 玩家

	PlayerEnter(player define.IPlayer)
	PlayerLeave(uid int64)
	OnClientSceneLoadOver(uid int64)

	// entity

	OnAddEntity(e entity.IEntity)
	OnDestroyEntity(e entity.IEntity)
	PushSnapshot(camera define.ICamera, e entity.IEntity)

	// 流程

	Start()
	IsOver() bool
	Update()
	Destroy()

	OnPreCameraEnter(camera define.ICamera)
	OnPostCameraEnter(camera define.ICamera)
}

type Creator func() ISceneLogic

type ISceneLogicFactory interface {
	Register(name string, c Creator)
	Create(name string) ISceneLogic
}
