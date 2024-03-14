package goid

import (
	"fmt"
	"testing"

	"github.com/faiface/mainthread"
	"github.com/stretchr/testify/assert"
)

func TestGoroutineId(t *testing.T) {
	mainThreadId := Get()
	fmt.Printf("main 协程id : %v\n", mainThreadId)

	for i := 0; i < 10; i++ {
		mainthread.Run(func() {
			mainthread.Call(func() {
				id := Get()
				fmt.Printf("mainThread 协程id : %v\n", id)

				assert.Equal(t, mainThreadId, id, "主协程ID不一致")
			})
		})
	}
}
