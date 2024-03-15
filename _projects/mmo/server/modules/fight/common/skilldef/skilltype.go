package skilldef

// 定义常量转化

const (
	SkillTypeInvalid = iota
	SkillTypeNormal
	SkillTypeTargetBullet
)

var (
	SkillTypeNames = []string{
		"",
		"normal",
		"targetbullet"}
)
