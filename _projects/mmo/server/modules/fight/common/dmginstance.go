package common

// DmgInstance 结算过程数据
// 结算先创建一个这个实例
// 存放各种临时数据
// 此结构会通过事件传递，各个插件可以基于逻辑修改数据
// 最终得到Dmg
// 还有各种flag,比如，暴击啊，之类
// 此结构在结算调用后会被立刻释放
type DmgInstance struct {
	Src      ICharacter
	SrcProxy ICharProxy
	Tar      ICharacter
	TarProxy ICharProxy

	// 计算数据
	// 根据技能配置，伤害类型
	// 采用不同的攻防计算
	DmgType int

	BonusPhysicPower float64 // 额外攻击
	BonusPenetrate   float64 // 物穿
	BonusPenetrateP  float64 // 额外物穿百分比
	TarBonusArmor    float64 // 额外物防

	BonusMagicPower      float64 // 额外法伤
	BonusMagicPenetrate  float64 // 额外法穿
	BonusMagicPenetrateP float64 // 额外法穿百分比
	TarBonusMagicArmor   float64 // 额外法防

	BonusCriticalRate      float64 // 额外暴击
	BonusCriticalDmgRate   float64 // 额外暴击伤害倍率
	BonusCriticalDmgResist float64 // 额外暴击伤害抵抗(对抗属性)

	BonusDmgEnhance    float64 // 额外增伤
	BonusDmgEnhanceP   float64 // 额外增伤百分比
	TarBonusDmgReduce  float64 // 额外减伤
	TarBonusDmgReduceP float64 // 额外减伤百分比

	// 结果
	NoDmg    bool // 实际未产生伤害
	Critical bool
	Dmg      int
	// 被吸收的伤害
	Absorb int

	TarBeKilled bool
}

func (i *DmgInstance) Reset() {
	i.Src = nil
	i.Tar = nil
	i.SrcProxy = nil
	i.TarProxy = nil
	i.DmgType = 0
	i.BonusPhysicPower = 0
	i.BonusPenetrate = 0
	i.BonusPenetrateP = 0
	i.TarBonusArmor = 0

	i.BonusMagicPower = 0
	i.BonusMagicPenetrate = 0
	i.BonusMagicPenetrateP = 0
	i.TarBonusMagicArmor = 0

	i.BonusCriticalRate = 0
	i.BonusCriticalDmgRate = 0
	i.BonusCriticalDmgResist = 0

	i.BonusDmgEnhance = 0
	i.BonusDmgEnhanceP = 0
	i.TarBonusDmgReduce = 0
	i.TarBonusDmgReduceP = 0

	i.NoDmg = false
	i.Critical = false
	i.Dmg = 0
	i.Absorb = 0
	i.TarBeKilled = false
}
