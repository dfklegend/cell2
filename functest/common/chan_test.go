package common

import (
	"log"
	"testing"
	"time"
)

// 测试chan
func Test_ChanMaxSize(t *testing.T) {
	c := make(chan int, 10)

	go func() {
		for i := 0; i < 13; i++ {
			log.Printf("in ")
			c <- 1
		}
		log.Printf("all in over ")
	}()

	go func() {
		for true {
			time.Sleep(time.Second)

			_ = <-c
		}
	}()

	time.Sleep(5 * time.Second)
}

//
