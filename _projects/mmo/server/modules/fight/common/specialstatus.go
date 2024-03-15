package common

const (
	SSNoSkill        int = iota // 禁止释放技能
	SSNoNormalAttack            // 禁止普通攻击
	SSNoMove                    // 禁止移动
	SSInvincible                // 无敌，不受伤害，免疫敌对效果

	SSMax // max
)

// ISpecialStatusCtrl impl in character/impls
type ISpecialStatusCtrl interface {
	AddSpecialStatus(id int)
	SubSpecialStatus(id int)
	HasSpecialStatus(id int) bool
}
