package skill

import (
	"mmo/modules/fight/common"

	"mmo/modules/csv"
)

func createSubSkill(owner common.ICharacter, parent *Skill, skillId string, level int, depth int) *Skill {
	cfg := csv.Skill.GetEntry(skillId)
	if cfg == nil {
		return nil
	}
	skill := NewSkill(skillId, level, cfg, owner)
	skill.SetTar(parent.tar)
	skill.SetWorld(owner.GetWorld())
	skill.SetSubSkill(true)
	return skill
}
