package common

// 测试map 安全性
// 结论: 竞争读会不会异常

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// 会错误退出
func Tes_t_CorReadWrite(t *testing.T) {

	m := make(map[int]int)

	fmt.Printf("%v\n", m[0])

	running := true
	temp := 0
	for i := 0; i < 100; i++ {
		go func() {
			for running {
				index := rand.Intn(100)
				if _, ok := m[index]; !ok {
					m[index] = 1
					continue
				}
				m[index]++
			}
		}()

		go func() {
			for running {
				index1 := rand.Intn(100)
				index2 := rand.Intn(100)
				if m[index1] < m[index2] {
					temp++
				}
			}
		}()
	}

	time.Sleep(10 * time.Second)
	running = false
}

// 只读取不会发生问题
func Test_CorRead(t *testing.T) {

	m := make(map[int]int)

	fmt.Printf("%v\n", m[0])

	running := true
	temp := 0

	for i := 0; i < 100; i++ {
		m[i] = i
	}

	for i := 0; i < 100; i++ {
		go func() {
			for running {
				index1 := rand.Intn(100)
				index2 := rand.Intn(100)
				if m[index1] < m[index2] {
					temp++
				}
			}
		}()
	}

	time.Sleep(10 * time.Second)
	running = false
}
