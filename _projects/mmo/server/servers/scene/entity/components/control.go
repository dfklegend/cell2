package components

import (
	define3 "mmo/servers/scene/define"
	"mmo/servers/scene/entity/define"
)

// ControlComponent
// 接受玩家的控制
type ControlComponent struct {
	*BaseSceneComponent
	move define.IMoveComponent
}

func NewControlComponent() *ControlComponent {
	return &ControlComponent{
		BaseSceneComponent: NewBaseSceneComponent(),
	}
}
func (c *ControlComponent) OnPrepare() {
	c.BaseSceneComponent.OnPrepare()
}

func (c *ControlComponent) OnStart() {
	c.move = c.GetOwner().GetComponent(define.MoveComponent).(define.IMoveComponent)
}

func (c *ControlComponent) ReqMoveTo(x, z float32) {
	c.move.MoveTo(define3.Pos{
		X: x,
		Z: z,
	})
}

func (c *ControlComponent) ReqStopMove(x, z float32) {
	// 应该存在偏差
	c.move.StopMove()
}
