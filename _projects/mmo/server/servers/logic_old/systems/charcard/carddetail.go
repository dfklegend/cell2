package charcard

//
//import (
//	"github.com/dfklegend/cell2/utils/event/light"
//	"github.com/dfklegend/cell2/utils/serialize/proto"
//
//	"mmo/common/cardattr"
//	"mmo/common/entity"
//	"mmo/common/fightutils"
//	"mmo/messages/cproto"
//	"mmo/modules/fight/attr"
//	"mmo/modules/fight/buf"
//	charimpls "mmo/modules/fight/character/impls"
//	"mmo/modules/fight/common"
//	"mmo/modules/fight/equip"
//	"mmo/modules/fight/skill"
//)
//
//const (
//	MaxEquip        = 6
//	MaxSkill        = 6
//	MainWeaponIndex = 0 // 0号装备是主武器
//)
//
//// CardDetail
//// 计算实际添加的属性
//type CardDetail struct {
//	owner *CharCard
//
//	data *cproto.CharCard
//	card common.ICharacter
//	dirt bool
//
//	events        *light.EventCenter
//	world         entity.IWorld
//	timeProvider  common.ITimeProvider
//	worldForFight common.IWorld
//
//	attrCWatcher *fightutils.AttrChangeWatcher
//}
//
//func newDetail() *CardDetail {
//	return &CardDetail{
//		dirt: false,
//	}
//}
//
//func (c *CardDetail) Init(owner *CharCard) {
//	c.owner = owner
//
//	c.events = light.NewEventCenter()
//	c.world = entity.NewWorld()
//	c.timeProvider = fightutils.NewTimeProvider()
//
//	c.worldForFight = fightutils.NewWorldForFight(c.world, c.timeProvider,
//		nil, nil, nil, nil)
//	c.attrCWatcher = fightutils.NewAttrWatcher()
//}
//
//func (c *CardDetail) setDirt() {
//	c.dirt = true
//}
//
//func (c *CardDetail) clearDirt() {
//	c.dirt = false
//}
//
//func (c *CardDetail) isDirt() bool {
//	return c.dirt
//}
//
//func (c *CardDetail) cloneData(from *cproto.CharCard) *cproto.CharCard {
//	serializer := proto.GetDefaultSerializer()
//	bytes, err := serializer.Marshal(from)
//	if err != nil {
//		return nil
//	}
//
//	data := &cproto.CharCard{}
//	serializer.Unmarshal(bytes, data)
//	ValidCard(data)
//	return data
//}
//
//func (c *CardDetail) saveTo(to *cproto.CharCard) bool {
//	if c.data == nil || !c.isDirt() {
//		return false
//	}
//
//	serializer := proto.GetDefaultSerializer()
//	bytes, err := serializer.Marshal(c.data)
//	if err != nil {
//		return false
//	}
//
//	serializer.Unmarshal(bytes, to)
//	c.clearDirt()
//	return true
//}
//
//func (c *CardDetail) createCharacter(data *cproto.CharCard) common.ICharacter {
//	b := charimpls.NewBuilder()
//	b.WithSkill(skill.NewCtrl())
//	b.WithSkillTable(skill.NewSkillTable())
//	b.WithBufCtrl(buf.NewCtrl())
//	b.WithSlots(equip.NewSlots(10))
//	card := b.Build()
//
//	card.Init(0, c.worldForFight, c.events)
//	cardattr.ApplyCard(card, data)
//
//	card.Start()
//	return card
//}
//
//func (c *CardDetail) Create(data *cproto.CharCard) {
//	c.data = c.cloneData(data)
//	c.card = c.createCharacter(data)
//	c.pushInitAttrs()
//}
//
//func (c *CardDetail) IsCreated() bool {
//	return c.card != nil
//}
//
//func (c *CardDetail) pushInitAttrs() {
//	if c.card == nil {
//		return
//	}
//
//	attrs := make([]*cproto.CardAttr, 0)
//	// 每一项属性
//	c.card.VisitAttr(func(index int, attr attr.IAttr) {
//		item := &cproto.CardAttr{
//			Index: int32(index),
//			Value: float32(attr.GetValue()),
//		}
//		attrs = append(attrs, item)
//	})
//
//	c.owner.PushCmd("initattrs", &cproto.CardAttrs{
//		Attrs: attrs,
//	})
//}
//
//func (c *CardDetail) SetEquip(index int, id string) {
//	if index < 0 || index >= MaxEquip {
//		return
//	}
//	if c.data.Equips == nil {
//		c.data.Equips = make([]*cproto.EquipSlot, MaxEquip)
//	}
//
//	item := c.data.Equips[index]
//	if item == nil {
//		item = &cproto.EquipSlot{}
//		c.data.Equips[index] = item
//	}
//
//	item.EquipId = id
//
//	c.prepareAttrWatcher()
//	if id != "" {
//		c.card.SetEquip(index, id)
//	} else {
//		c.card.RemoveEquip(index)
//	}
//	c.pushAttrAndClearWatcher()
//
//	c.setDirt()
//}
//
//func (c *CardDetail) pushAttrChanged(watcher *fightutils.AttrChangeWatcher) {
//	attrs := make([]*cproto.CardAttr, 0)
//	watcher.Visit(func(one *fightutils.Attr) {
//		item := &cproto.CardAttr{
//			Index: int32(one.Index),
//			Value: float32(one.NewV),
//		}
//		attrs = append(attrs, item)
//	})
//
//	if len(attrs) == 0 {
//		return
//	}
//	c.owner.PushCmd("attrs", &cproto.CardAttrs{
//		Attrs: attrs,
//	})
//}
//
//func (c *CardDetail) prepareAttrWatcher() {
//	c.attrCWatcher.Reset()
//	c.card.AddAttrChangeWatcher(c.attrCWatcher)
//}
//
//func (c *CardDetail) pushAttrAndClearWatcher() {
//	c.pushAttrChanged(c.attrCWatcher)
//	c.card.RemoveAttrChangeWatcher(c.attrCWatcher)
//}
//
//// SetSkill	技能没有属性变化
//func (c *CardDetail) SetSkill(index int, id string, level int32) bool {
//	if index < 0 || index >= MaxSkill {
//		return false
//	}
//
//	data := c.data
//	if data.Skills == nil {
//		data.Skills = make([]*cproto.SkillSlot, MaxSkill)
//	}
//
//	// 检查是否有同一个技能多次添加
//	if id != "" {
//		for _, v := range data.Skills {
//			if v != nil && v.SkillId == id {
//				return false
//			}
//		}
//	}
//
//	item := data.Skills[index]
//	if item == nil {
//		item = &cproto.SkillSlot{}
//		data.Skills[index] = item
//	}
//
//	item.SkillId = id
//	if id == "" {
//		item.Level = 0
//	} else {
//		item.Level = level
//	}
//	c.setDirt()
//	return true
//}
