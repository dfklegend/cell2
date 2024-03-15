package common

// IFightDetailRecorder
// 在战斗内负责记录详细战报
// 战斗结束后，统一输出
// 便于诊断
type IFightDetailRecorder interface {
	OnStartFight()

	OnStartSkill(id SkillId, src ICharacter, tar ICharacter)
	OnSkillHit(id SkillId, src ICharacter, tar ICharacter, dmg int)

	OnPreAddBuf(c ICharacter, id BufId, level, stack int)
	// OnPostAddBuf 输出属性变化
	OnPostAddBuf(c ICharacter, id BufId, level, stack int)

	OnEndFight()
	Visit(doFunc func(info string))
}
