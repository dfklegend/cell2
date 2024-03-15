package frame

// IFrameClock
// 驱动固定帧
type IFrameClock interface {
	Start(fixedFrameTime int64)
	Update(nowMs int64)
	TryForwardFrame() bool

	GetCurFrameTime() int64
	GetCurFrameIndex() int32
	NowFrameMs() int64
}

// IFrameLogic
// 被FrameDriver驱动
type IFrameLogic interface {
	GetCurFrameTime() int64
	GetCurFrameIndex() int32
	NowMs() int64
}
