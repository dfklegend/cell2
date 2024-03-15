package states

import (
	ai2 "mmo/servers/scene/entity/ai"
	"mmo/servers/scene/entity/ai/ctrl"
	"mmo/servers/scene/entity/components"
	define2 "mmo/servers/scene/entity/define"
	fsm2 "mmo/servers/scene/entity/fsm"
)

func init() {
	ai2.GetStateFactory().Register(ai2.StateNoMoveAttack, fsm2.NewFuncCreator(func() fsm2.IState {
		return NewNoMoveAttack()
	}))
}

// 根据攻击间隔，判断是否攻击
// 移动时不能攻击

type NoMoveAttack struct {
	*fsm2.BaseState

	aiCtrl    *ctrl.AICtrl
	nextCheck int64
}

func NewNoMoveAttack() fsm2.IState {
	return &NoMoveAttack{
		BaseState: fsm2.NewBaseState(),
	}
}

func (a *NoMoveAttack) OnEnter() {
	a.aiCtrl = a.GetCtrl().GetContext().(*ctrl.AICtrl)
}

func (a *NoMoveAttack) OnLeave() {
	ai := a.aiCtrl
	ai.SetEnemy(-1)
	ai.GetUnit().ClearTar()
	// break now attack
}

func (a *NoMoveAttack) Update() {
	// 看看有没有敌人
	if a.IsOver() {
		return
	}

	if a.keepInAttackRange() {
		a.tryAttack()
	}
}

func (a *NoMoveAttack) IsOver() bool {
	ai := a.aiCtrl
	if ai.GetUnit().IsDead() {
		return true
	}
	if ai.GetMove().IsMoving() {
		return true
	}

	enemyId := ai.GetEnemy()
	if enemyId == 0 {
		return true
	}
	owner := ai.GetOwner()
	enemy := owner.GetWorld().GetEntity(enemyId)
	if enemy == nil {
		return true
	}
	eUnit := ai.GetEnemyUnit()
	if eUnit.IsDead() {
		return true
	}
	return false
}

// true: in attack range
func (a *NoMoveAttack) keepInAttackRange() bool {
	ai := a.aiCtrl
	enemyId := ai.GetEnemy()
	owner := ai.GetOwner()
	enemy := owner.GetWorld().GetEntity(enemyId)
	if enemy == nil {
		return false
	}

	pos := ai.GetTran().GetPos()
	enemyPos := enemy.GetComponent(define2.Transform).(*components.Transform).GetPos()

	dist := pos.Distance(enemyPos)
	if dist <= a.aiCtrl.GetAttackRange() {
		return true
	}
	return false
}

func (a *NoMoveAttack) tryAttack() {
	ai := a.aiCtrl
	if ai.GetMove().IsMoving() {
		return
	}
	ai.GetUnit().UpdateAttack(ai.GetEnemy())
}
