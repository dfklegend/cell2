package impl

import (
	"mmo/common/entity"
)

type BaseComponent struct {
	entity.IComponent
	owner entity.IEntity
}

func NewBaseComponent() *BaseComponent {
	return &BaseComponent{}
}

func (c *BaseComponent) SetOwner(e entity.IEntity) {
	c.owner = e
}

func (c *BaseComponent) GetOwner() entity.IEntity {
	return c.owner
}

func (c *BaseComponent) OnPrepare() {
}

func (c *BaseComponent) OnStart() {
}

func (c *BaseComponent) OnDestroy() {
}

func (c *BaseComponent) Update() {
}

func (c *BaseComponent) LateUpdate() {
}
