package impl

import (
	"mmo/common/entity"
)

type Entity struct {
	entity.IEntity
	world     entity.IWorld
	id        entity.EntityID
	destroyed bool

	components map[string]entity.IComponent

	quickComps []entity.IComponent
	dirt       bool

	started bool
}

func NewEntity() *Entity {
	return &Entity{
		destroyed:  false,
		components: map[string]entity.IComponent{},
		dirt:       false,
		started:    false,
	}
}

func (e *Entity) SetWorld(world entity.IWorld) {
	e.world = world
}

func (e *Entity) GetWorld() entity.IWorld {
	return e.world
}

func (e *Entity) SetId(id entity.EntityID) {
	e.id = id
}

func (e *Entity) GetId() entity.EntityID {
	return e.id
}

func (e *Entity) IsDestroyed() bool {
	return e.destroyed
}

func (e *Entity) AddComponent(name string, component entity.IComponent) entity.IComponent {
	component.SetOwner(e)
	e.components[name] = component

	e.dirt = true

	if e.isStarted() {
		// 后面增加的component
		component.OnPrepare()
		component.OnStart()
	}
	return component
}

// FilterComponents TODO: 是否需要保证component顺序
func (e *Entity) FilterComponents(do func(c entity.IComponent)) {
	if e.dirt {
		e.quickComps = make([]entity.IComponent, 0)
		for _, v := range e.components {
			e.quickComps = append(e.quickComps, v)
		}
		e.dirt = false
	}

	for i := 0; i < len(e.quickComps); i++ {
		do(e.quickComps[i])
	}
}

func (e *Entity) Prepare() {
	e.FilterComponents(func(c entity.IComponent) {
		c.OnPrepare()
	})
}

func (e *Entity) Start() {
	e.FilterComponents(func(c entity.IComponent) {
		c.OnStart()
	})
	e.started = true
}

func (e *Entity) isStarted() bool {
	return e.started
}

func (e *Entity) Update() {
	e.FilterComponents(func(c entity.IComponent) {
		c.Update()
	})
}

func (e *Entity) LateUpdate() {
	e.FilterComponents(func(c entity.IComponent) {
		c.LateUpdate()
	})
}

func (e *Entity) DestroySelf() {
	e.GetWorld().DestroyEntity(e.GetId())
}

// DoDestroy only can call by world
func (e *Entity) DoDestroy() {
	if e.destroyed {
		return
	}

	e.FilterComponents(func(c entity.IComponent) {
		c.OnDestroy()
	})
	e.components = map[string]entity.IComponent{}
	e.quickComps = nil

	e.destroyed = true
}

func (e *Entity) GetComponent(name string) entity.IComponent {
	return e.components[name]
}
