package components

import (
	"mmo/servers/scene/entity/ai"
	define2 "mmo/servers/scene/entity/define"
)

type AIComponent struct {
	*BaseSceneComponent
	aiCtrl ai.IAICtrl
}

func NewAIComponent() *AIComponent {
	return &AIComponent{
		BaseSceneComponent: NewBaseSceneComponent(),
	}
}

func (c *AIComponent) InitAICtrl(ctrl ai.IAICtrl) {
	c.aiCtrl = ctrl
	ctrl.InitCtrl(c.GetOwner())
}

func (c *AIComponent) ActiveCtrl(b bool) {
	c.aiCtrl.SetActive(b)
}

func (c *AIComponent) OnPrepare() {
	c.BaseSceneComponent.OnPrepare()
	c.aiCtrl.OnPrepare()

	unit := c.GetOwner().GetComponent(define2.BaseUnit).(*BaseUnit)

	unit.events.SubscribeWithReceiver("onDead", c, c.onDead)
	unit.events.SubscribeWithReceiver("onRelive", c, c.onRelive)
}

func (c *AIComponent) OnStart() {
	c.aiCtrl.OnStart()
}

func (c *AIComponent) Update() {
	c.aiCtrl.Update()
}

func (c *AIComponent) OnDestroy() {
}

func (c *AIComponent) onDead(args ...any) {

}

func (c *AIComponent) onRelive(args ...any) {

}
