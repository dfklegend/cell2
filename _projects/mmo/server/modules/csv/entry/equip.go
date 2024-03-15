package entry

import (
	"mmo/modules/csv/structs"
)

type Equip struct {
	Id   string `csv:"id"`
	Name string `csv:"name"`

	Buf string `csv:"buf"`

	Attr0 structs.AttrValue `csv:"attr0"`
	V0    float32           `csv:"v0"`

	Attr1 structs.AttrValue `csv:"attr1"`
	V1    float32           `csv:"v1"`

	Attr2 structs.AttrValue `csv:"attr2"`
	V2    float32           `csv:"v2"`

	Attr3 structs.AttrValue `csv:"attr3"`
	V3    float32           `csv:"v3"`
}

func (b *Equip) GetId() string {
	return b.Id
}
