package benchtest

import (
	"math/rand"
	"testing"
)

// 测试闭包额外的性能损失有多少，有没有额外内存消耗
// 结论
//     并没有额外的分配，局部的调用，在栈上

func doFunc1(i int, op func() int) {
	oldV := i
	newV := op()
	onGetValue(oldV, newV)
}

func doFunc(i int) {
	oldV := i
	i++
	i -= rand.Intn(100)
	newV := i
	onGetValue(oldV, newV)
}

func onGetValue(oldV, newV int) {
	// malloc in stack
	v := make([]int, 100)
	v[0] = 100
}

func Benchmark_Closure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := i
		doFunc1(v, func() int {
			v++
			v -= rand.Intn(100)
			return v
		})
	}
}

func Benchmark_Normal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doFunc(i)
	}
}
