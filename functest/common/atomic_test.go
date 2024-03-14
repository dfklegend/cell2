package common

import (
	"log"
	"sync/atomic"
	"testing"
	"time"
)

func Test_Add(t *testing.T) {
	var v int32
	for i := 0; i < 100000; i++ {
		go func() {
			v++
		}()
	}
	time.Sleep(time.Second)
	log.Println(atomic.LoadInt32(&v))
}

func Test_AtomicAdd(t *testing.T) {
	var v int32
	for i := 0; i < 100000; i++ {
		go func() {
			atomic.AddInt32(&v, 1)
		}()
	}
	time.Sleep(time.Second)
	log.Println(atomic.LoadInt32(&v))
}
