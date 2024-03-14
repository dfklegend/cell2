package proto

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dfklegend/cell2/utils/serialize/proto/msgs"
)

func TestNewSerializer(t *testing.T) {
	t.Parallel()

	log.Println(msgs.ServiceRequest{})

	serializer := NewSerializer()
	assert.NotNil(t, serializer)
}

func TestNormal(t *testing.T) {
	in := &msgs.ServiceRequest{
		Sender: "someone",
	}

	serializer := NewSerializer()
	bytes, _ := serializer.Marshal(in)

	out := &msgs.ServiceRequest{}
	serializer.Unmarshal(bytes, out)
	assert.Equal(t, "someone", out.Sender)
}

// 测试不同版本的结构定义会如何
// 测试结果:
//	序列化数据不匹配，不会异常
// 	前面类型一致，会读出来
// 	序号和类型匹配的会被读取出来
// 	类型不一致，数据会被重置成初始化状态
func TestMismatch(t *testing.T) {
	in := &msgs.ServiceRequest{
		Sender: "someone",
		ReqId:  99,
		Type:   "type",
	}

	log.Printf("origin data: %+v\n", in)

	serializer := NewSerializer()
	bytes, _ := serializer.Marshal(in)

	r1 := &msgs.ServiceRequest1{}
	serializer.Unmarshal(bytes, r1)
	log.Printf("%+v\n", r1)

	assert.Equal(t, "someone", r1.Sender)
	assert.Equal(t, "", r1.Type)

	r2 := &msgs.ServiceRequest2{}
	r2.ReqId = 100
	serializer.Unmarshal(bytes, r2)
	log.Printf("%+v\n", r2)

	// 类型不一致，会被设成默认值
	assert.Equal(t, int32(0), r2.ReqId)
	// 类型和序号一致，能读取出来
	assert.Equal(t, "type", r2.Type)
}

// 122 ns/op
func BenchmarkMarshal(b *testing.B) {
	in := &msgs.ServiceRequest{
		Sender: "someone",
		ReqId:  99,
		Type:   "type",
	}

	serializer := NewSerializer()
	for i := 0; i < b.N; i++ {
		serializer.Marshal(in)
	}
}

// 149 ns/op
func BenchmarkUnMarshal(b *testing.B) {
	in := &msgs.ServiceRequest{
		Sender: "someone",
		ReqId:  99,
		Type:   "type",
	}

	serializer := NewSerializer()
	ret, _ := serializer.Marshal(in)
	out := &msgs.ServiceRequest{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		serializer.Unmarshal(ret, out)
	}
}
