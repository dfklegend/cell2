package entry

import (
	"mmo/modules/csv/base"
	"mmo/modules/csv/structs/skilleffect"
)

type SkillEffect struct {
	*base.IArgs
	SkillId   string                `csv:"skillId"`
	ApplyTime skilleffect.ApplyTime `csv:"applyTime"`
	TarType   skilleffect.TarType   `csv:"tarType"`
	Op        skilleffect.OpType    `csv:"op"`
}

func (s *SkillEffect) GetId() string {
	return ""
}
