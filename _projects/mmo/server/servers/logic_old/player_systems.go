package logic_old

import (
	"mmo/messages/cproto"
	"mmo/modules/fight/common"
	"mmo/servers/logic_old/systems/bridge"
)

func (p *Player) makeSystemInterfaces() {
	// 初始化快捷接口
	//p.charCard = p.systems.GetSystem("charcard").(bridge.ICharCard)
}

func (p *Player) GetCharCard() bridge.ICharCard {
	return p.charCard
}

func (p *Player) CreateCharCard(name string) int32 {
	return p.charCard.CreateCard(name)
}

func (p *Player) DeleteCharCard(id int32) bool {
	return p.charCard.DeleteCard(id)
}

func (p *Player) OpenCharCard(id int32) {
	p.charCard.OpenCard(id)
}

func (p *Player) GetCharCardNum() int32 {
	return p.charCard.GetCardNum()
}

func (p *Player) GetCharCardBrief() string {
	return p.charCard.Brief()
}

func (p *Player) GetCardByIndex(index int) *cproto.CharCard {
	return p.charCard.GetCardByIndex(index)
}

func (p *Player) GetCard(id int32) *cproto.CharCard {
	return p.charCard.GetCard(id)
}

func (p *Player) CardSetEquip(index int, id common.EquipId) {
	p.charCard.SetEquip(index, id)
}

func (p *Player) CardSave() {
	p.charCard.SaveCard()
}

func (p *Player) SystemRequest(system, cmd string, args []byte, cb func(ret []byte, errCode int32)) {
	//p.systems.Request(system, cmd, args, cb)
}
