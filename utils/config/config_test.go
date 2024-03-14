package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	config := LoadConfig("../data/config/")
	fmt.Printf("%v,%v,%v\n", config.GetInt("exam.int"),
		config.GetString("exam.int"), config.GetString("exam.string"))

	assert.Equal(t, 10, config.GetInt("exam.int"))
	// 转换成了string
	assert.Equal(t, "10", config.GetString("exam.int"))

	fmt.Printf("%v\n", config.GetInt("exam.int2", 99))
	assert.Equal(t, 99, config.GetInt("exam.int2", 99))

	fmt.Printf("%v\n", config.GetInt("exam.int2"))
	// 不存在
	assert.Equal(t, 0, config.GetInt("exam.int2"))

	fmt.Printf("%v\n", config.GetString("exam.int2", "def"))
	assert.Equal(t, "def", config.GetString("exam.int2", "def"))
}
