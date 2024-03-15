package functest

import (
	"testing"
)

// 可以做一些功能测试

func TestSome(t *testing.T) {
	m := map[string]any{}
	m["a"] = 1
	b := m["b"]
	print(m["a"])
	print(b)
	//assert.Equal(t, 90, 90)
}
