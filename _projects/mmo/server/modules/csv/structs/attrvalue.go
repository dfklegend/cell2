package structs

import (
	"strings"

	"mmo/modules/fight/common"
)

type AttrValue struct {
	Type      int
	IsPercent bool
}

func (a *AttrValue) getIndex(attrName string) int {
	return common.AttrNameToIndex(attrName)
}

//
//	生命:  某个属性， 生命,%  百分比属性

func (a *AttrValue) UnmarshalCSV(csv string) (err error) {
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
