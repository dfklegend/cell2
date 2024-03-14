package common

import (
	"encoding/json"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetStackStr() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

func SafeJsonMarshalByteArray(data interface{}) []byte {
	r, e := json.Marshal(data)
	if e != nil {
		return []byte("")
	}
	return r
}

func SafeJsonMarshal(data interface{}) string {
	r, e := json.Marshal(data)
	if e != nil {
		return ""
	}
	return string(r)
}

func NowMs() int64 {
	return time.Now().UnixNano() / 1e6
}

//	NowNano 纳秒
func NowNano() int64 {
	return time.Now().UnixNano()
}

func RandFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func SplitAddress(address string) (string, int) {
	subs := strings.Split(address, ":")
	if len(subs) != 2 {
		return "", 0
	}
	port, _ := strconv.Atoi(subs[1])
	return subs[0], port
}
