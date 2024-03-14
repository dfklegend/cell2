package convert

import (
	"strconv"
)

func TryParseInt64(s string, def int64) int64 {
	ret, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return ret
}

func TryParseInt32(s string, def int32) int32 {
	ret, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return def
	}
	return int32(ret)
}

func TryParseInt(s string, def int) int {
	ret, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return def
	}
	return int(ret)
}
