package entry

import (
	"mmo/modules/csv/structs"
)

type Skill struct {
	Id           string                `csv:"id"`
	Name         string                `csv:"name"`
	SubSkills    structs.SubSkillsType `csv:"subSkills"`
	CD           float32               `csv:"cd"`
	Type         structs.SkillType     `csv:"type"`
	TarType      structs.TarType       `csv:"tarType"`
	NormalAttack bool                  `csv:"normalAttack"`
	TotalTime    int32                 `csv:"totalTime"`
	HitTime      int32                 `csv:"hitTime"`
	DmgType      structs.DmgType       `csv:"dmgType"`
	BaseDmg      float32               `csv:"baseDmg"`
	BaseDmgLv    float32               `csv:"baseDmgLv"`
	AD           float32               `csv:"ad"`
	ADLv         float32               `csv:"adLv"`
	AP           float32               `csv:"ap"`
	APLv         float32               `csv:"apLv"`
}

func (s *Skill) GetId() string {
	return s.Id
}
