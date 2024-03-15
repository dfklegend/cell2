package frame

type FrameLock struct {
	lastUpdateTime int64
	localTime      int64

	interval   int64
	frameTime  int64
	frameIndex int64
}

func (l *FrameLock) Start(interval int64, nowMs int64) {
	l.interval = interval
	l.frameTime = 0
	l.frameIndex = 0

	l.lastUpdateTime = nowMs
	l.localTime = 0
}

func (l *FrameLock) Update(nowMs int64) {
	offset := nowMs - l.lastUpdateTime
	if offset < 0 {
		return
	}

	l.localTime += offset
}

func (l *FrameLock) TryForwardFrame() bool {
	offset := l.localTime - l.frameTime
	if offset < l.interval {
		return false
	}

	l.frameIndex++
	l.frameTime += l.interval
	return true
}
