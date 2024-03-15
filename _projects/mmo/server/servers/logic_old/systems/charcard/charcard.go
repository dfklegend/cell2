package charcard

//
//import (
//	"fmt"
//	"reflect"
//	"strings"
//
//	"github.com/dfklegend/cell2/utils/event/light"
//	"golang.org/x/exp/slices"
//
//	"mmo/common/define"
//	"mmo/common/factory"
//	mymsg "mmo/messages"
//	"mmo/messages/cproto"
//	"mmo/modules/fight/common"
//	"mmo/servers/logic_old/systems"
//	sysfactory "mmo/servers/logic_old/systems/factory"
//)
//
//func init() {
//	sysfactory.GetFactory().RegisterFunc("charcard", func(args ...any) factory.IObject {
//		return newCharCard("charcard")
//	})
//}
//
//func Visit() {}
//
//// CharCard
//// 角色卡牌的数据存储
//type CharCard struct {
//	*systems.BaseSystem
//
//	nextCardId int32
//	cards      []*cproto.CharCard
//
//	cardOpened int32
//	detail     *CardDetail
//}
//
//func newCharCard(name string) *CharCard {
//	return &CharCard{
//		BaseSystem: systems.NewBaseSystem(name),
//		detail:     newDetail(),
//	}
//}
//
//func (c *CharCard) OnCreate() {
//	c.detail.Init(c)
//
//	c.bindEvents(true)
//
//	// cmds
//	c.RegisterCmdHandler("", reflect.TypeOf(&cproto.TestAdd{}), c.onCmdAdd)
//	c.RegisterCmdHandler("opencard", reflect.TypeOf(&cproto.TestAdd{}), c.onCmdAdd)
//	c.RegisterCmdHandler("setequip", reflect.TypeOf(&cproto.CardSetEquip{}), c.onCmdSetEquip)
//	c.RegisterCmdHandler("setskill", reflect.TypeOf(&cproto.CardSetSkill{}), c.onCmdSetSkill)
//	c.RegisterCmdHandler("savecard", reflect.TypeOf(&cproto.EmptyArg{}), c.onCmdSaveCard)
//}
//
//func (c *CharCard) OnDestroy() {
//}
//
//func (c *CharCard) bindEvents(bind bool) {
//	events := c.GetPlayer().GetEvents()
//	light.BindEventWithReceiver(bind, events, "onEnterWorld", c, c.onEventEnterWorld)
//}
//
//func (c *CharCard) PushInfoToClient() {
//}
//
//func (c *CharCard) OnEnterWorld() {
//
//}
//
//func (c *CharCard) onEventEnterWorld(args ...any) {
//}
//
//func (c *CharCard) LoadData(info *mymsg.PlayerInfo) {
//	c.cards = info.Cards
//	if c.cards == nil {
//		c.cards = make([]*cproto.CharCard, 0)
//	}
//	c.nextCardId = info.NextCardId
//	if c.nextCardId == 0 {
//		c.nextCardId = 1
//	}
//
//	c.validCards()
//}
//
//func (c *CharCard) validCards() {
//	// 保证数据结构有效
//	for i := 0; i < len(c.cards); i++ {
//		card := c.cards[i]
//		ValidCard(card)
//	}
//}
//
//func (c *CharCard) SaveData(info *mymsg.PlayerInfo) {
//	if len(c.cards) > 0 {
//		info.Cards = c.cards
//	}
//	if c.nextCardId > 1 {
//		info.NextCardId = c.nextCardId
//	}
//}
//
//func (c *CharCard) allocId() int32 {
//	id := c.nextCardId
//	c.nextCardId++
//	return id
//}
//
//func (c *CharCard) GetCards() []*cproto.CharCard {
//	return c.cards
//}
//
//func (c *CharCard) GetCardByIndex(index int) *cproto.CharCard {
//	if index < 0 || index >= len(c.cards) {
//		return nil
//	}
//	return c.cards[index]
//}
//
//func (c *CharCard) GetCard(id int32) *cproto.CharCard {
//	index := c.findIndex(id)
//	if index == -1 {
//		return nil
//	}
//	return c.cards[index]
//}
//
//func (c *CharCard) CreateCard(name string) int32 {
//	card := &cproto.CharCard{
//		Id:    c.allocId(),
//		Name:  name,
//		Level: 1,
//		Stat:  &cproto.StatInfo{},
//	}
//	ValidCard(card)
//	c.cards = append(c.cards, card)
//	c.GetPlayer().SetDirt()
//
//	c.RefreshCards()
//	return card.Id
//}
//
//func (c *CharCard) findIndex(id int32) int {
//	return slices.IndexFunc(c.cards, func(c *cproto.CharCard) bool {
//		return c.Id == id
//	})
//}
//
//func (c *CharCard) DeleteCard(id int32) bool {
//	index := c.findIndex(id)
//	if index == -1 {
//		return false
//	}
//
//	c.cards = slices.Delete(c.cards, index, index+1)
//	c.GetPlayer().SetDirt()
//
//	c.RefreshCards()
//	return true
//}
//
//func (c *CharCard) GetCardNum() int32 {
//	return int32(len(c.cards))
//}
//
//func (c *CharCard) Brief() string {
//	var builder strings.Builder
//	builder.Grow(1000)
//	builder.WriteString(fmt.Sprintf("num %v\n", len(c.cards)))
//	for _, v := range c.cards {
//		builder.WriteString(fmt.Sprintf("  Id: %v Name: %v\n", v.Id, v.Name))
//	}
//	return builder.String()
//}
//
//func (c *CharCard) RefreshCards() {
//	p := c.GetPlayer()
//	p.PushMsg("refreshcards", &cproto.RefreshCards{
//		Cards: c.GetCards(),
//	})
//}
//
//func (c *CharCard) UpdateAndRefresh() {
//	c.GetPlayer().SetDirt()
//	c.RefreshCards()
//}
//
//func (c *CharCard) onCmdAdd(srcArgs any, cb func(ret any, code int32)) {
//}
//
//// open card to edit
//func (c *CharCard) onCmdOpenCard(srcArgs any, cb func(ret any, code int32)) {
//
//}
//
//func (c *CharCard) OpenCard(id int32) {
//	data := c.GetCard(id)
//	if data == nil {
//		return
//	}
//	c.cardOpened = id
//	c.detail.Create(data)
//}
//
//func (c *CharCard) onCmdSetEquip(srcArgs any, cb func(ret any, code int32)) {
//	args, _ := srcArgs.(*cproto.CardSetEquip)
//	code := define.Succ
//	if !c.SetEquip(int(args.Index), args.Id) {
//		code = define.ErrFaild
//	}
//	cb(&cproto.EmptyArg{}, int32(code))
//}
//
//func (c *CharCard) SetEquip(index int, id common.EquipId) bool {
//	if !c.detail.IsCreated() {
//		return false
//	}
//	c.detail.SetEquip(index, id)
//	return true
//}
//
//func (c *CharCard) onCmdSetSkill(srcArgs any, cb func(ret any, code int32)) {
//	args, _ := srcArgs.(*cproto.CardSetSkill)
//	code := define.Succ
//
//	if !c.SetSkill(int(args.Index), args.Id, args.Level) {
//		code = define.ErrFaild
//	}
//
//	cb(&cproto.EmptyArg{}, int32(code))
//}
//
//func (c *CharCard) SetSkill(index int, id common.SkillId, level int32) bool {
//	if !c.detail.IsCreated() {
//		return false
//	}
//	return c.detail.SetSkill(index, id, level)
//}
//
//func (c *CharCard) onCmdSaveCard(srcArgs any, cb func(ret any, code int32)) {
//	c.SaveCard()
//	cb(&cproto.EmptyArg{}, int32(define.Succ))
//}
//
//func (c *CharCard) SaveCard() {
//	data := c.GetCard(c.cardOpened)
//	if data == nil {
//		return
//	}
//	if c.detail.saveTo(data) {
//		c.GetPlayer().SetDirt()
//	}
//}
