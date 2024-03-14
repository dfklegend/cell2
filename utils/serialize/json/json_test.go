// Copyright (c) nano Author and TFG Co. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package json

import (
	"encoding/json"
	"log"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dfklegend/cell2/utils/serialize/proto/msgs"
)

func TestNewSerializer(t *testing.T) {
	t.Parallel()

	serializer := NewSerializer()

	assert.NotNil(t, serializer)
}

func TestMarshal(t *testing.T) {
	t.Parallel()

	type MyStruct struct {
		Str    string
		Number float64
	}
	var marshalTables = map[string]struct {
		raw       interface{}
		marshaled []byte
		errType   interface{}
	}{
		"test_ok": {
			&MyStruct{Str: "hello", Number: 42},
			[]byte(`{"Str":"hello","Number":42}`),
			nil,
		},
		"test_nok": {
			&MyStruct{Number: math.Inf(1)},
			nil,
			&json.UnsupportedValueError{},
		},
	}
	serializer := NewSerializer()

	for name, table := range marshalTables {
		t.Run(name, func(t *testing.T) {
			result, err := serializer.Marshal(table.raw)

			assert.Equal(t, table.marshaled, result)
			if table.errType == nil {
				assert.NoError(t, err)
			} else {
				assert.IsType(t, table.errType, err)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	type MyStruct struct {
		Str    string
		Number int
	}
	var unmarshalTables = map[string]struct {
		data        []byte
		unmarshaled *MyStruct
		errType     interface{}
	}{
		"test_ok": {
			[]byte(`{"Str":"hello","Number":42}`),
			&MyStruct{Str: "hello", Number: 42},
			nil,
		},
		"test_nok": {
			[]byte(`invalid`),
			nil,
			&json.SyntaxError{},
		},
	}
	serializer := NewSerializer()

	for name, table := range unmarshalTables {
		t.Run(name, func(t *testing.T) {
			var result MyStruct
			err := serializer.Unmarshal(table.data, &result)
			if table.errType == nil {
				assert.NoError(t, err)
				assert.Equal(t, table.unmarshaled, &result)
			} else {
				assert.Empty(t, &result)
				assert.IsType(t, table.errType, err)
			}
		})
	}
}

type ServiceRequest struct {
	Sender string
	ReqId  int32
	Type   string
}

type ServiceRequest1 struct {
	ReqId  int32
	Sender string
	Type   string
}

type ServiceRequest2 struct {
	ReqId  string
	Sender string
	Type   int32
}

/*
	测试数据结构变化, json序列化接口的影响
	名字一致，读取正确
*/
func TestMismatch(t *testing.T) {
	in := &ServiceRequest{
		Sender: "someone",
		ReqId:  99,
		Type:   "type",
	}

	log.Printf("origin data: %+v\n", in)

	serializer := NewSerializer()
	bytes, _ := serializer.Marshal(in)

	r1 := &ServiceRequest1{}
	serializer.Unmarshal(bytes, r1)
	log.Printf("%+v\n", r1)

	assert.Equal(t, "someone", r1.Sender)
	assert.Equal(t, "type", r1.Type)

	r2 := &ServiceRequest2{}
	r2.ReqId = "100"
	serializer.Unmarshal(bytes, r2)
	log.Printf("%+v\n", r2)

	// 类型不一致，会被忽略
	assert.Equal(t, "100", r2.ReqId)
	assert.Equal(t, int32(0), r2.Type)
}

// 223 ns/op
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
