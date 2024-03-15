package common

// ISkillTable
// 可用技能列表
// 分析当前可用技能，选择下一个技能
type ISkillTable interface {
	Init(owner ICharacter, provider ITimeProvider)

	SetNormalAttackSkill(id SkillId)
	SetNormalAttackInterval(interval float32)

	AddSkill(id SkillId, level int)
	UpgradeSkill(id SkillId, level int)
	RemoveSkill(id SkillId)

	Update()
	GetNextSkill() (skillId SkillId, level int)

	// skill cd

	PushCD(id SkillId, prefireTime int32)
	GetCDRest(id SkillId) float32
	GetCDRestPercent(id SkillId) float32
	OffsetCD(id SkillId, offset float32)
	OffsetCDPercent(id SkillId, offset float32)
	IsCDReady(skillId SkillId) bool
}
