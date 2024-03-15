package common

/*
	如果直接用common.NowMs/1000.0会有巨大的浮点数精度问题
	由于浮点数精度分布，在大数值情况下精度损失很大

	不要用浮点数来计时
*/
type ITimeProvider interface {
	Update()
	NowMs() int64
}
