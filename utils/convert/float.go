package convert

import (
	"strconv"
)

func TryParseFloat64(s string, def float64) float64 {
	ret, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def
	}
	return ret
}

func TryParseFloat32(s string, def float32) float32 {
	ret, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return def
	}
	return float32(ret)
}
