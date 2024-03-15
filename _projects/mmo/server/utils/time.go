package utils

import (
	"time"
)

func GetAbsoluteDay() int32 {
	return int32(GetAbsoluteTime() / 86400)
}

func GetAbsoluteTime() int64 {
	return int64(GetTime() / 1000)
}

func GetTime() int64 {
	return int64(time.Now().UnixNano() / 1000000)
}
