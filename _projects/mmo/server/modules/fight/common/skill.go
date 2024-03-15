package common

type ISkillCtrl interface {
	Update()

	//GetSkillCDRestTime(id SkillId) int

	IsSkillRunning() bool
	// CastSkill 释放技能
	CastSkill(id SkillId, level int, src ICharacter, tar CharId)
	CallbackSkill(id SkillId, level int, src ICharacter, tar CharId)
	BreakSkill(src ICharacter, breakNormalAttack bool)
	//CastSkillAtPos(id SkillId, level int, tarPos Pos)
	// BackgroundSkill 后台技能
	//BackgroundSkill(id SkillId)
}
