package charimpls

import (
	"mmo/modules/fight/common"
)

func (c *Character) IsInvincible() bool {
	return c.HasSpecialStatus(common.SSInvincible)
}

func (c *Character) OnSpecialStatusChanged(id int, oldV bool, newV bool) {
	switch id {
	case common.SSNoSkill:
		if newV {
			c.onGotNoSkill()
		}
	case common.SSNoNormalAttack:
		if newV {
			c.onGotNoNormalAttack()
		}
	}
}

func (c *Character) onGotNoSkill() {
	// 取消技能
	c.BreakSkill(c, false)
}

func (c *Character) onGotNoNormalAttack() {
	// 取消普攻
	c.BreakSkill(c, true)
}
