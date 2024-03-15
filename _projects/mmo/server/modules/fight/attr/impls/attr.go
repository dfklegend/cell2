package impls

import (
	"math"

	"mmo/modules/fight/attr"
)

// 考虑:
// 属性是否需要不同实现
// 比如，有的属性只是整数字段，也不会%缩放

type Attr struct {
	base    attr.Value
	percent float64
	dirt    bool

	valueClamper   attr.IValueClamper
	percentClamper attr.IPercentClamper
	watcher        attr.IWatcher

	final attr.Value
}

func NewAttr() *Attr {
	return &Attr{
		dirt: false,
	}
}

func (a *Attr) Copy(v1 *Attr) *Attr {
	v1.base = a.base
	v1.percent = a.percent
	v1.dirt = a.dirt
	v1.valueClamper = a.valueClamper
	v1.percentClamper = a.percentClamper
	v1.watcher = a.watcher
	v1.final = a.final
	return v1
}

func (a *Attr) Reset() {
	a.base = 0
	a.percent = 0
	a.setDirt()
}

func (a *Attr) SetPercentClamper(clamper attr.IPercentClamper) {
	a.percentClamper = clamper
	a.setDirt()
}

func (a *Attr) SetValueClamper(clamper attr.IValueClamper) {
	a.valueClamper = clamper
	a.setDirt()
}

func (a *Attr) SetWatcher(watcher attr.IWatcher) {
	a.watcher = watcher
}

func (a *Attr) SetBase(v attr.Value) {
	a.base = v
	a.setDirt()
}

func (a *Attr) setDirt() {
	a.dirt = true
	if a.watcher != nil {
		a.watcher.OnDirt(a)
	}
}

func (a *Attr) OffsetBase(v attr.Value) {
	a.base += v
	a.setDirt()
}

func (a *Attr) SetPercent(v float32) {
	a.percent = float64(v)
	a.setDirt()
}

func (a *Attr) OffsetPercent(v float32) {
	a.percent += float64(v)
	a.setDirt()
}

func (a *Attr) GetValue() attr.Value {
	if a.dirt {
		a.calcFinal()
		a.dirt = false
	}

	return a.final
}

func (a *Attr) GetIntValue() int {
	return int(math.Ceil(a.GetValue() - 0.5))
}

func (a *Attr) calcFinal() {
	percent := 1 + a.percent
	if a.percentClamper != nil {
		percent = a.percentClamper.Clamp(percent)
	}
	// 逻辑意义上，最少衰减到0%
	if percent < 0 {
		percent = 0
	}

	final := a.base * percent
	if a.valueClamper != nil {
		final = a.valueClamper.Clamp(final)
	}

	a.final = final
}

func (a *Attr) IsDirt() bool {
	return a.dirt
}
