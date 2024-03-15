package ctrl

import (
	"github.com/dfklegend/cell2/utils/common"

	"mmo/servers/scene/entity/ai"
	"mmo/servers/scene/entity/fsm"
	"mmo/servers/scene/space/utils"
)

func init() {

}

// states
// StateInit, CardStateWait, CardStateAttack

func CardNormal(ctrl fsm.ICtrl) {
	aiCtrl := ctrl.GetContext().(*AICtrl)
	state := ctrl.GetStateType()
	now := common.NowMs()
	if state == ai.StateInit {
		ctrl.ChangeState(ai.CardStateWait)
		return
	}

	if state == ai.CardStateWait {
		// 根据间隔，进行索敌
		if now >= aiCtrl.Options.NextCanSearch {
			aiCtrl.Options.NextCanSearch = now + aiCtrl.Options.SearchInterval
			enemyId := utils.FindNearestEnemy(aiCtrl.GetSpace(), aiCtrl.GetOwner(), 999)
			if enemyId > 0 {
				aiCtrl.SetEnemy(enemyId)
				ctrl.ChangeState(ai.CardStateAttack)
			}
		}
		return
	}

	if state == ai.CardStateAttack {
		if ctrl.GetState().IsOver() {
			if !aiCtrl.IsDead() {
				ctrl.ChangeState(ai.CardStateWait)
			}
		}
		return
	}
}
