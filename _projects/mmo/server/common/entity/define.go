package entity

// 提供一个基本的Entity实现
// 便于对象组合
// 后续可以考虑提供component pool
// 还有 component的update, 避免空转

type EntityID int32

const (
	InvalidEntityId = 0
)

// 不要动态添加component
// 而是统一加好，然后prepare

type IEntity interface {
	SetWorld(world IWorld)
	GetWorld() IWorld

	SetId(id EntityID)
	GetId() EntityID
	IsDestroyed() bool

	AddComponent(name string, component IComponent) IComponent
	Prepare()
	Start()

	Update()
	LateUpdate()

	DestroySelf()
	DoDestroy()

	GetComponent(name string) IComponent
}

type IComponent interface {
	SetOwner(e IEntity)
	GetOwner() IEntity

	// OnPrepare 做一些数据初始化
	OnPrepare()
	// OnStart 启动，可以获得entity所有的component
	// 建议在这里获取关注的其他component
	OnStart()
	Update()
	LateUpdate()
	OnDestroy()
}

// IWorldContext world的使用context，外部传入
type IWorldContext interface {
}

type IWorld interface {
	SetContext(ctx IWorldContext)
	GetContext() IWorldContext

	AllocId() EntityID

	SetEventListener(listener IWorldEventListener)

	AddEntity(e IEntity)
	GetEntity(id EntityID) IEntity
	DestroyEntity(id EntityID)

	Update()
	Destroy()

	Filter(doFunc func(entity IEntity))
}

type IWorldEventListener interface {
	OnAddEntity(entity IEntity)
	OnDestroyEntity(entity IEntity)
}
