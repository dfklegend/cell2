package skilleffect

// ----
// 效果目标
const (
	TarTypeInvalid = iota
	TarCaster
	TarTarget
)

var (
	TarTypeNames = []string{
		"",
		"caster",
		"target",
	}
)

// ----
// 触发时机
const (
	ATInvalid      = iota
	CastCheck      // 施放时触发
	HitCasterCheck // 命中触发一次(aoe也只触发一次)
	HitTargetCheck // 命中每个目标触发(aoe可能多个目标)
	OnFailedCheck  // 技能失败时触发
)

var (
	ApplyTimeNames = []string{
		"",
		"cast_check",
		"hit_caster_check",
		"hit_target_check",
		"onfailed_check",
	}
)
