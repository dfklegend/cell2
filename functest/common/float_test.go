package common

import (
	"log"
	"testing"
)

func Test_Float(t *testing.T) {
	i64 := 1669359089646
	f64 := float64(i64) / 1000
	f32 := float32(f64)

	f32_1 := float32(i64) / 1000

	log.Printf("%v %f %f %f\n", i64, f64, f32, f32_1)
	// 1669359089646 1669359089.646000 1669359104.000000
	// 大浮点数损失精度超乎想象

	// 所以如果要用time.Now().UnixNano来计时，float32会有问题
	// 要么改进为小数值范围，比如进程启动时间等

	// func NowMs() int64 {
	//	return time.Now().UnixNano() / 1e6
	//}
}
