package define

import (
	"mmo/modules/csv/entry"
)

type ISkill interface {
	GetCfg() *entry.Skill
}
