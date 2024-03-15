package fightutils

import (
	"github.com/dfklegend/cell2/utils/common"
)

//

type TimeProvider struct {
	beginTime int64
	now       int64
}

func NewTimeProvider() *TimeProvider {
	return &TimeProvider{
		beginTime: common.NowMs(),
	}
}

func (p *TimeProvider) Update() {
	p.now = p.getPassed()
}

func (p *TimeProvider) NowMs() int64 {
	return p.now
}

func (p *TimeProvider) getPassed() int64 {
	return common.NowMs() - p.beginTime
}

/*
	如果直接用common.NowMs/1000.0会有巨大的浮点数精度问题
	由于浮点数精度分布，在大数值情况下精度损失很大
*/
/*func (p *TimeProvider) Now() float32 {
	now := p.NowMs()
	v := float64(now) / 1000.0
	v32 := float32(v)
	return v32
}*/
