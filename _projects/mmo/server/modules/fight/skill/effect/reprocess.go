package effect

import (
	"mmo/modules/csv"
	"mmo/modules/csv/base"
	"mmo/modules/csv/entry"
	"mmo/modules/fight/skill/effect/factory"
)

func ProcessArgs() {
	csv.SkillEffect.Visit(func(entry *entry.SkillEffect) {
		op := factory.GetOpFactory().Create(entry.Op.Type)
		if op != nil {
			argsFormatter, ok := op.(base.IArgsFormatter)
			if ok {
				argsFormatter.Format(entry.IArgs)
			}
		}
	})
}
