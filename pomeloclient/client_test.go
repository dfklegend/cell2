package client

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/dfklegend/cell2/pomelonet/common/conn/message"
	"github.com/dfklegend/cell2/utils/serialize/json"
)

type QueryGateReq struct {
}

type QueryGateAck struct {
	Code int    `json:"code"`
	IP   string `json:"ip"`
	Port string `json:"port"`
}

// 开启一下_examples\chat2
// 功能测试
func TestRawClient(t *testing.T) {
	c := NewCellClient("client", json.GetDefaultSerializer())
	go func() {
		c.Start("127.0.0.1:30021")
		c.WaitReady()

		c.GetClient().SendRequest("gate.gate.querygate", []byte("{}"), func(error bool, msg *message.Message) {
			fmt.Println("ack from cb")
			fmt.Println(string(msg.Data))
		})
	}()

	time.Sleep(5 * time.Second)
	c.Stop()
}

// 测试
func TestCellClient(t *testing.T) {
	c := NewCellClient("client", json.GetDefaultSerializer())
	go func() {
		c.Start("127.0.0.1:30021")
		c.WaitReady()

		c.SendRequest("gate.gate.querygate", &QueryGateReq{}, func(err error, ret interface{}) {
			fmt.Println("ack from cb")
			ack := ret.(*QueryGateAck)
			str, _ := json.GetDefaultSerializer().Marshal(ack)
			fmt.Println(string(str))
		}, reflect.TypeOf(&QueryGateAck{}))
	}()

	time.Sleep(5 * time.Second)
	c.Stop()
}
