package states

import (
	"mmo/servers/scene/define"
	ai2 "mmo/servers/scene/entity/ai"
	"mmo/servers/scene/entity/ai/ctrl"
	"mmo/servers/scene/entity/components"
	define2 "mmo/servers/scene/entity/define"
	fsm2 "mmo/servers/scene/entity/fsm"
)

func init() {
	ai2.GetStateFactory().Register(ai2.StateAttack, fsm2.NewFuncCreator(func() fsm2.IState {
		return NewAttack()
	}))
}

// 根据攻击间隔，判断是否攻击
// 如果不在攻击范围，移动过去

type Attack struct {
	*fsm2.BaseState

	aiCtrl    *ctrl.AICtrl
	nextCheck int64
}

func NewAttack() fsm2.IState {
	return &Attack{
		BaseState: fsm2.NewBaseState(),
	}
}

func (a *Attack) OnEnter() {
	a.aiCtrl = a.GetCtrl().GetContext().(*ctrl.AICtrl)
}

func (a *Attack) OnLeave() {
	ai := a.aiCtrl
	ai.SetEnemy(-1)
	ai.GetUnit().ClearTar()
}

func (a *Attack) Update() {
	// 看看有没有敌人
	if a.IsOver() {
		return
	}

	if a.keepInAttackRange() {
		a.tryAttack()
	}
}

func (a *Attack) IsOver() bool {
	ai := a.aiCtrl
	if ai.GetUnit().IsDead() {
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
func (a *Attack) keepInAttackRange() bool {
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
		if dist <= 0.3*ai.GetAttackRange() {
			// 太近了，离远点
			// 最佳距离 0.8的attackrange
			a.moveToBestPos(enemyPos)
			return false
		}

		if ai.GetMove().IsMoving() {
			ai.GetMove().StopMove()
		}
		return true
	}

	moveC := ai.GetMove()
	moving := moveC.IsMoving()
	if moving {
		return false
	}

	moveC.MoveTo(enemyPos)
	return false
}

func (a *Attack) moveToBestPos(enemyPos define.Pos) {
	pos := a.aiCtrl.GetTran().GetPos()
	dir := enemyPos.Sub(pos)
	dir = dir.Normalized()

	bestPos := pos.Add(dir.Mul(a.aiCtrl.GetAttackRange() * 0.8))
	a.aiCtrl.GetMove().MoveTo(bestPos)
}

func (a *Attack) tryAttack() {
	ai := a.GetCtrl().GetContext().(*ctrl.AICtrl)
	//if ai.GetUnit().IsSkillRunning() {
	//	return
	//}
	//
	////unit.Hit(eUnit)
	//ai.GetUnit().StartSkill(1, ai.GetEnemy())
	ai.GetUnit().UpdateAttack(ai.GetEnemy())
}
