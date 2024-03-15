package ctrl

import (
	"github.com/dfklegend/cell2/utils/common"

	"mmo/servers/scene/entity/ai"
	"mmo/servers/scene/entity/fsm"
	"mmo/servers/scene/space/utils"
)

func init() {

}

// Player
// 不会主动移动
// wait, nomoveattack, dead
func Player(ctrl fsm.ICtrl) {
	aiCtrl := ctrl.GetContext().(*AICtrl)
	state := ctrl.GetStateType()
	now := common.NowMs()
	if state == ai.StateInit {
		ctrl.ChangeState(ai.StateWait)
		return
	}

	if state == ai.StateWait {
		// 根据间隔，进行索敌
		if now >= aiCtrl.Options.NextCanSearch {
			aiCtrl.Options.NextCanSearch = now + aiCtrl.Options.SearchInterval
			enemyId := utils.FindNearestEnemy(aiCtrl.GetSpace(), aiCtrl.GetOwner(), aiCtrl.GetAttackRange())
			if enemyId > 0 {
				aiCtrl.SetEnemy(enemyId)
				ctrl.ChangeState(ai.StateNoMoveAttack)
			}
		}
		return
	}

	if state == ai.StateNoMoveAttack {
		if ctrl.GetState().IsOver() {
			if aiCtrl.IsDead() {
				ctrl.ChangeState(ai.StateDead)
			} else {
				ctrl.ChangeState(ai.StateWait)
			}
		}
		return
	}

	if state == ai.StateDead {
		if !aiCtrl.IsDead() {
			ctrl.ChangeState(ai.StateWait)
		}
	}
}
