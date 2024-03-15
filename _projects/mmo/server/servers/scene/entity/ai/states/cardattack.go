package states

import (
	ai2 "mmo/servers/scene/entity/ai"
	"mmo/servers/scene/entity/ai/ctrl"
	fsm2 "mmo/servers/scene/entity/fsm"
)

func init() {
	ai2.GetStateFactory().Register(ai2.CardStateAttack, fsm2.NewFuncCreator(func() fsm2.IState {
		return NewCardAttack()
	}))
}

// 根据攻击间隔，判断是否攻击
// 如果不在攻击范围，移动过去

type CardAttack struct {
	*fsm2.BaseState

	aiCtrl *ctrl.AICtrl
}

func NewCardAttack() fsm2.IState {
	return &CardAttack{
		BaseState: fsm2.NewBaseState(),
	}
}

func (a *CardAttack) OnEnter() {
	a.aiCtrl = a.GetCtrl().GetContext().(*ctrl.AICtrl)
}

func (a *CardAttack) OnLeave() {
	ai := a.aiCtrl
	ai.SetEnemy(-1)
	ai.GetUnit().ClearTar()
}

func (a *CardAttack) Update() {
	// 看看有没有敌人
	if a.IsOver() {
		return
	}

	a.tryAttack()
}

func (a *CardAttack) IsOver() bool {
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

func (a *CardAttack) tryAttack() {
	ai := a.aiCtrl
	//if ai.GetUnit().IsSkillRunning() {
	//	return
	//}
	//
	////unit.Hit(eUnit)
	//ai.GetUnit().StartSkill(1, ai.GetEnemy())
	ai.GetUnit().UpdateAttack(ai.GetEnemy())
}
