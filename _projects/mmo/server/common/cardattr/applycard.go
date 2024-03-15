package cardattr

import (
	"mmo/messages/cproto"
	"mmo/modules/fight/attrinitor"
	"mmo/modules/fight/common"
)

func ApplyCard(char common.ICharacter, card *cproto.CharCard) {
	char.ChangeLevel(1)
	char.AddEquipGroup(attrinitor.NewBaseModelInitor())

	// equips
	if card.Equips != nil && len(card.Equips) > 0 {
		for i, v := range card.Equips {
			if v.EquipId != "" {
				char.SetEquip(i, v.EquipId)
			}
		}
	}

	// skills
	skills := char.GetSkillTable()

	if card.NormalSkill != nil && card.NormalSkill.SkillId != "" {
		skills.SetNormalAttackSkill(card.NormalSkill.SkillId)
	} else {
		skills.SetNormalAttackSkill("普攻")
	}

	if card.Skills != nil && len(card.Skills) > 0 {
		for _, v := range card.Skills {
			if v.SkillId != "" {
				skills.AddSkill(v.SkillId, int(v.Level))
			}
		}
	}
	//skills.AddSkill("英勇打击", 1)
	//skills.AddSkill("嗜血", 1)
	//skills.AddSkill("致命打击", 1)
}
