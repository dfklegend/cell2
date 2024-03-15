package common

// 属性定义
const (
	Level int = iota
	Side
	HP
	HPMax
	Energy
	EnergyMax

	AttackSpeed  // 攻速 基础值受武器影响
	CriticalRate // 暴击概率 [0-1]

	PhysicPower // 物伤
	MagicPower  // 法伤
	PhysicArmor
	MagicArmor

	WeaponMinDmg // 武器最小伤害
	WeaponMaxDmg // 武器最大伤害

	MaxAttrNum // last
)

func init() {
	registerAttrs()
}

func registerAttrs() {
	var R = RegisterAttr
	R(Level, "等级")
	R(Side, "阵营")
	R(HP, "-")
	// 属性定义中一般血量都是hpmax
	R(HPMax, "生命")
	R(Energy, "-")
	R(EnergyMax, "能量")

	R(AttackSpeed, "攻速")
	R(CriticalRate, "暴击")

	R(PhysicPower, "物伤")
	R(MagicPower, "法伤")
	R(PhysicArmor, "护甲")
	R(MagicArmor, "法抗")

	R(WeaponMinDmg, "武器最小伤害")
	R(WeaponMaxDmg, "武器最大伤害")
}
