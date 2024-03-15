package charcard

import (
	"mmo/messages/cproto"
)

func ValidCard(card *cproto.CharCard) {
	// 保证数据结构有效
	if card.Stat == nil {
		card.Stat = &cproto.StatInfo{}
	}
	if card.NormalSkill == nil {
		card.NormalSkill = &cproto.SkillSlot{}
	}
}
