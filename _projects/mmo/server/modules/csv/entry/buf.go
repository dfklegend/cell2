package entry

import (
	"mmo/modules/csv/structs"
)

type Buf struct {
	Id             string                  `csv:"id"`
	Name           string                  `csv:"name"`
	MaxStack       int32                   `csv:"maxStack"`
	DeadKeep       bool                    `csv:"deadKeep"`
	Times          int                     `csv:"times"`
	Interval       int                     `csv:"interval"`
	AffectType     int                     `csv:"affectType"`
	TriggleSkillId string                  `csv:"triggleSkillId"`
	SkillTar       structs.BufSkillTarType `csv:"skillTar"`

	SpecialStatus structs.IntListType `csv:"specialstatus"`
	Attr0         structs.AttrValue   `csv:"attr0"`
	Base0         float32             `csv:"base0"`
	Lv0           float32             `csv:"lv0"`

	Attr1 structs.AttrValue `csv:"attr1"`
	Base1 float32           `csv:"base1"`
	Lv1   float32           `csv:"lv1"`

	Attr2 structs.AttrValue `csv:"attr2"`
	Base2 float32           `csv:"base2"`
	Lv2   float32           `csv:"lv2"`
}

func (b *Buf) GetId() string {
	return b.Id
}
