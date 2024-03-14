package csvutils

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 可以自定义数据类型，只要提供UnmarshalCSV接口
// 如果数据不能转化，比如类型是int，数据是一个无法转换的string，会panic

type Hero struct {
	Id    int    `csv:"id"`
	Name  string `csv:"name"`
	Model string `csv:"model"`
}

type Equip struct {
	Id       string  `csv:"id"`
	EType    int     `csv:"eType"`
	BaseType string  `csv:"baseType"`
	Speed    float32 `csv:"speed"`
	MinDmg   int     `csv:"minDmg"`
	MaxDmg   int     `csv:"maxDmg"`
	Block    int     `csv:"block"`
	Attr0    Attr    `csv:"attr0"`
	Value0   float32 `csv:"value0"`
	Attr1    Attr    `csv:"attr1"`
	Value1   float32 `csv:"value1"`
	Attr2    string  `csv:"attr2"`
	Value2   float32 `csv:"value2"`
	Attr3    string  `csv:"attr3"`
	Value3   float32 `csv:"value3"`
}

// 数据转化
type Attr struct {
	Type      int
	IsPercent bool
}

func (a *Attr) getIndex(csv string) int {
	switch csv {
	case "耐力":
		return 1
	case "力量":
		return 2
	case "敏捷":
		return 3
	}
	return 0
}

func (a *Attr) UnmarshalCSV(csv string) (err error) {
	subs := strings.Split(csv, ",")
	if len(subs) < 1 {
		return nil
	}

	a.Type = a.getIndex(subs[0])
	if len(subs) > 1 {
		a.IsPercent = subs[1] == "%"
	}
	return nil
}

func TestReadEquipCSV(t *testing.T) {
	var equips []*Equip
	err := LoadFromFile("equip.csv", &equips)

	if err != nil {
		panic(err)
	}

	for i := range equips {
		equip := equips[i]
		log.Println(equip)
	}

	assert.Equal(t, 4, len(equips))
}

func TestReadTestCSV(t *testing.T) {

	var heroes []*Hero
	err := LoadFromFileWithIndex("test.csv", &heroes, 2, 3)

	if err != nil {
		panic(err)
	}
	for _, hero := range heroes {
		log.Println(hero)
	}

	assert.Equal(t, 1, len(heroes))
	assert.Equal(t, "hero_102", heroes[0].Model)
}
