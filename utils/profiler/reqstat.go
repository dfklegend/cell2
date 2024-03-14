package profiler

import (
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/logger"
)

//	ReqStat 统计调用情况
type ReqStat struct {
	beginTime int64
	endTime   int64

	// 	总共调用次数
	called int64
	// 	总共调用时间
	totalResponseTime int64
}

func NewReqStat() *ReqStat {
	return &ReqStat{}
}

func (s *ReqStat) Begin() {
	s.beginTime = common.NowMs()
	s.endTime = 0
	s.called = 0
	s.totalResponseTime = 0
}

func (s *ReqStat) AddCall(cost int64) {
	s.called++
	s.totalResponseTime += cost
}

func (s *ReqStat) End() {
	s.endTime = common.NowMs()
}

func (s *ReqStat) Report(name string) {
	l := logger.Log
	if s.called == 0 {
		l.Infof(" no stat")
		return
	}

	l.Infof("----- %v req stat -------", name)
	l.Infof("called :%v time :%v", s.called, s.endTime-s.beginTime)
	l.Infof("total response time :%v", s.totalResponseTime)
	l.Infof("avg response time :%.2f", float64(s.totalResponseTime)/float64(s.called))
	l.Infof("------------")
}
