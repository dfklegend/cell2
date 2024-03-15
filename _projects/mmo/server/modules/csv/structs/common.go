package structs

import (
	"strings"

	"github.com/dfklegend/cell2/utils/convert"
)

// StringListType 字符串列表
type StringListType struct {
	Strs []string
}

func (t *StringListType) UnmarshalCSV(csv string) (err error) {
	if strings.TrimSpace(csv) == "" {
		return nil
	}
	t.Strs = strings.Split(csv, ",")
	return nil
}

// IntListType 字符串列表
type IntListType struct {
	Ints []int
}

func (t *IntListType) UnmarshalCSV(csv string) (err error) {
	if strings.TrimSpace(csv) == "" {
		return nil
	}
	strs := strings.Split(csv, ",")
	t.Ints = make([]int, len(strs))
	for i := 0; i < len(strs); i++ {
		t.Ints[i] = convert.TryParseInt(strs[i], -1)
	}
	return nil
}
