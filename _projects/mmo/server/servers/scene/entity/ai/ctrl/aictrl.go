package ctrl

import (
	"mmo/common/entity"
	"mmo/servers/scene/define"
	ai2 "mmo/servers/scene/entity/ai"
	components2 "mmo/servers/scene/entity/components"
	define2 "mmo/servers/scene/entity/define"
	fsm2 "mmo/servers/scene/entity/fsm"
	"mmo/servers/scene/events"
)

func init() {

}

func NewAICtrl(f fsm2.CtrlFunc) ai2.IAICtrl {
	ctrl := &AICtrl{
		BaseCtrl: fsm2.NewBaseCtrl(),
		Options:  NewDefaultOptions(),
		active:   true,
	}
	ctrl.Init(ai2.GetStateFactory(), f, ctrl)
	return ctrl
}

// AICtrl 一些公共数据
type AICtrl struct {
	*fsm2.BaseCtrl

	owner entity.IEntity
	scene define.IScene
	space define.ISpace

	enemy entity.EntityID

	tran  *components2.Transform
	move  define2.IMoveComponent
	unit  *components2.BaseUnit
	skill *components2.SkillComponent

	eTran *components2.Transform
	eUnit *components2.BaseUnit

	Options *AIOptions

	active bool
}

func (c *AICtrl) SetActive(a bool) {
	c.active = a
}

func (c *AICtrl) InitCtrl(owner entity.IEntity) {
	c.owner = owner
}

func (c *AICtrl) GetOwner() entity.IEntity {
	return c.owner
}

func (c *AICtrl) GetSpace() define.ISpace {
	return c.space
}

func (c *AICtrl) GetEnemy() entity.EntityID {
	return c.enemy
}

func (c *AICtrl) SetEnemy(id entity.EntityID) {
	e := c.owner.GetWorld().GetEntity(id)
	if e == nil {
		return
	}
	c.enemy = id

	c.eTran = e.GetComponent(define2.Transform).(*components2.Transform)
	c.eUnit = e.GetComponent(define2.BaseUnit).(*components2.BaseUnit)
}

func (c *AICtrl) OnPrepare() {
}

func (c *AICtrl) OnStart() {
	e := c.owner
	if e.GetComponent(define2.MoveComponent) != nil {
		c.move = e.GetComponent(define2.MoveComponent).(define2.IMoveComponent)
	}

	c.tran = e.GetComponent(define2.Transform).(*components2.Transform)
	c.unit = e.GetComponent(define2.BaseUnit).(*components2.BaseUnit)
	c.skill = e.GetComponent(define2.Skill).(*components2.SkillComponent)

	tran := c.tran
	c.scene = tran.GetScene()
	c.space = c.scene.GetSpace()

	c.bindEvents()
}

func (c *AICtrl) bindEvents() {
	events := c.unit.GetEvents()

	events.SubscribeWithReceiver("behit", c, c.onBeHit)
	events.SubscribeWithReceiver("onDead", c, c.onDead)
	events.SubscribeWithReceiver("onRelive", c, c.onRelive)
}

func (c *AICtrl) Update() {
	if c.active {
		c.BaseCtrl.Update()
	}
}

func (c *AICtrl) GetTran() *components2.Transform {
	return c.tran
}

func (c *AICtrl) GetEnemyTran() *components2.Transform {
	return c.eTran
}

func (c *AICtrl) GetMove() define2.IMoveComponent {
	return c.move
}

func (c *AICtrl) GetUnit() *components2.BaseUnit {
	return c.unit
}

func (c *AICtrl) GetSkill() *components2.SkillComponent {
	return c.skill
}

func (c *AICtrl) GetEnemyUnit() *components2.BaseUnit {
	return c.eUnit
}

func (c *AICtrl) onBeHit(args ...any) {
	//l.L.Infof("onBeHit")
	e := args[0].(*events.EventBeHit)
	src := e.Src.(*components2.BaseUnit)
	if c.enemy == 0 {
		c.SetEnemy(src.GetOwner().GetId())
	}
}

func (c *AICtrl) onDead(args ...any) {
	c.ChangeState(ai2.StateDead)
}

func (c *AICtrl) onRelive(args ...any) {
	c.ChangeState(ai2.StateInit)
}

func (c *AICtrl) IsDead() bool {
	return c.unit.IsDead()
}

func (c *AICtrl) GetAttackRange() float32 {
	return 1.5
}
