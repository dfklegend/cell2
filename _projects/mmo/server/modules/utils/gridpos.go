package utils

import (
	"math"

	define3 "mmo/servers/scene/define"
)

func MakeGridValue(v float32) float32 {
	return float32(math.Round(float64(v)))
}

func MakeGridPos(pos define3.Pos) define3.Pos {
	return define3.Pos{
		X: MakeGridValue(pos.X),
		Z: MakeGridValue(pos.Z),
	}
}
