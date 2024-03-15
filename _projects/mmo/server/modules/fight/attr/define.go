package attr

type Value = float64

// IAttr 属性
// 可以百分比加成
type IAttr interface {
	Reset()

	// SetPercentClamper 比如, 百分比衰减最少减到 10%
	SetPercentClamper(clamper IPercentClamper)
	// SetValueClamper 比如, 最大血量最少也是1
	SetValueClamper(clamper IValueClamper)

	SetWatcher(watcher IWatcher)

	SetBase(v Value)
	OffsetBase(v Value)

	SetPercent(v float32)
	OffsetPercent(v float32)

	GetValue() Value
	GetIntValue() int

	IsDirt() bool
}

type IPercentClamper interface {
	Clamp(v float64) float64
}

type IValueClamper interface {
	Clamp(v Value) Value
}

// IWatcher
// 可以监视某个属性的变化
type IWatcher interface {
	OnDirt(attr IAttr)
}
