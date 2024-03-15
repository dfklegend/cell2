package skilldef

const (
	DmgTypeInvalid = iota
	DmgTypeNoDmg
	DmgTypePhysic
	DmgTypeMagic
	DmgTypeReal
)

var (
	DmgTypeNames = []string{
		"",
		"无伤害",
		"物理",
		"法术",
		"真伤",
	}
)
