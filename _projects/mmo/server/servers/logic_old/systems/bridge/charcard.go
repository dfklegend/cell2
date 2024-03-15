package bridge

import (
	"mmo/messages/cproto"
	"mmo/modules/fight/common"
)

type ICharCard interface {
	GetCardNum() int32
	CreateCard(name string) int32
	DeleteCard(id int32) bool
	OpenCard(id int32)
	Brief() string

	GetCards() []*cproto.CharCard
	GetCardByIndex(index int) *cproto.CharCard
	GetCard(id int32) *cproto.CharCard

	UpdateAndRefresh()

	SetEquip(index int, id common.EquipId) bool
	SaveCard()
}
