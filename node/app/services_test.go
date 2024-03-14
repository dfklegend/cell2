package app

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dfklegend/cell2/node/cluster"
	"github.com/dfklegend/cell2/nodectrl/define"
)

func makeMembers() []*cluster.Member {
	members := make([]*cluster.Member, 0)

	one := &cluster.Member{
		Id:    "node-1",
		Host:  "127.0.0.1",
		Port:  1024,
		State: int(define.Working),
		Services: []string{
			"chat.chat1",
			"chat.chat2",
			"gate.gate1",
			"gate.gate2",
		},
	}
	members = append(members, one)

	two := &cluster.Member{
		Id:   "node-2",
		Host: "127.0.0.1",
		Port: 1025,
		Services: []string{
			"chat.chat3",
			"gate.gate3",
			"gate.gate4",
		},
	}
	members = append(members, two)
	return members
}

func TestMake(t *testing.T) {

	members := makeMembers()
	m1, m2, m3, m4 := MakeMembers(members)

	log.Println(m1)
	log.Println(m2)
	log.Println(m3)
	log.Println(m4)

	assert.Equal(t, 3, len(m2["chat"].Items))
	assert.Equal(t, 2, len(m3["chat"].Items))
	assert.Equal(t, 4, len(m2["gate"].Items))
	assert.Equal(t, 2, len(m3["gate"].Items))
}

// 整体对象赋值，所以，并不存在锁的问题
func TestSync(t *testing.T) {
	mgr := NewClusterServices()
	members := makeMembers()

	mgr.MakeMembers(members)
	running := true

	for i := 0; i < 100; i++ {
		go func() {
			for running {
				chats := mgr.GetServiceList("chat")
				_ = len(chats.Items)
				//time.Sleep(1 * time.Millisecond)
			}
		}()

		go func() {
			for running {
				gates := mgr.GetServiceList("gate")
				_ = len(gates.Items)
				//time.Sleep(10 * time.Millisecond)
			}
		}()

		go func() {
			for running {
				mgr.MakeMembers(members)
				//time.Sleep(10 * time.Millisecond)
			}
		}()
	}

	time.Sleep(5 * time.Second)
	running = false
}
