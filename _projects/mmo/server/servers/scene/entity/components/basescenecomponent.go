package components

import (
	"mmo/common/entity/impl"
	"mmo/servers/scene/define"
)

type BaseSceneComponent struct {
	*impl.BaseComponent
	scene define.IScene
}

func NewBaseSceneComponent() *BaseSceneComponent {
	return &BaseSceneComponent{
		BaseComponent: impl.NewBaseComponent(),
	}
}

func (c *BaseSceneComponent) OnPrepare() {
	c.scene = c.GetOwner().GetWorld().GetContext().(define.IScene)
}

func (c *BaseSceneComponent) GetScene() define.IScene {
	return c.scene
}
