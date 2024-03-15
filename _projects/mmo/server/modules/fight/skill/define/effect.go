package define

import (
	"mmo/modules/csv/entry"
	"mmo/modules/fight/common"
)

// ISkillEffect 技能在不同时机可以触发一些效果
// 比如，添加buf
// 每种效果有自己的参数配置
type ISkillEffect interface {
	Apply(caster common.ICharacter, tar common.ICharacter, cfg *entry.SkillEffect, skillLv int)
}
