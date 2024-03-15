package ctrl

import (
	"github.com/dfklegend/cell2/utils/common"

	"mmo/servers/scene/entity/ai"
	"mmo/servers/scene/entity/fsm"
	"mmo/servers/scene/space/utils"
)

func init() {

}

func Normal(ctrl fsm.ICtrl) {
	aiCtrl := ctrl.GetContext().(*AICtrl)
	state := ctrl.GetStateType()
	now := common.NowMs()
	if state == ai.StateInit {
		ctrl.ChangeState(ai.StateRandMove)
		return
	}

	if state == ai.StateRandMove {
		// 根据间隔，进行索敌
		if now >= aiCtrl.Options.NextCanSearch {
			aiCtrl.Options.NextCanSearch = now + aiCtrl.Options.SearchInterval
			enemyId := utils.FindNearestEnemy(aiCtrl.GetSpace(), aiCtrl.GetOwner(), 6)
			if enemyId > 0 {
				aiCtrl.SetEnemy(enemyId)
				ctrl.ChangeState(ai.StateAttack)
			}
		}
		return
	}

	if state == ai.StateAttack {
		if ctrl.GetState().IsOver() {
			if !aiCtrl.IsDead() {
				ctrl.ChangeState(ai.StateRandMove)
			}
		}
		return
	}
}
