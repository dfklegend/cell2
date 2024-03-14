package mailbox

import (
	"time"
)

func NowMs() int64 {
	return time.Now().UnixNano() / 1e6
}

//	NowNano 纳秒
func NowNano() int64 {
	return time.Now().UnixNano()
}
