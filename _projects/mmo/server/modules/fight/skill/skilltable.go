package skill

import (
	"golang.org/x/exp/slices"

	"mmo/modules/csv"
	"mmo/modules/csv/entry"
	"mmo/modules/fight/common"
)

const (
	MaxSlotNum = 99
)

type Item struct {
	SkillId      common.SkillId
	Level        int
	Cfg          *entry.Skill
	CD           int64
	TimeCDFinish int64
}

func allocItem() *Item {
	return &Item{}
}

func freeItem(i *Item) {
}

// 	TODO: skillCD需要独立出来(增删不会影响CD)
type Table struct {
	owner    common.ICharacter
	provider common.ITimeProvider

	normalAttack *Item
	skills       []*Item

	nextSkill *Item

	normalAttackInterval int64
}

func NewSkillTable() *Table {
	return &Table{
		skills: []*Item{},
	}
}

func (t *Table) Init(owner common.ICharacter, provider common.ITimeProvider) {
	t.owner = owner
	t.provider = provider
}

func (t *Table) setTimeProvider(provider common.ITimeProvider) {
	t.provider = provider
}

func (t *Table) findCfg(id common.SkillId) *entry.Skill {
	return csv.Skill.GetEntry(id)
}

func (t *Table) initItem(item *Item, id common.SkillId, level int) {
	item.SkillId = id
	item.Level = level
	item.Cfg = t.findCfg(id)
	if item.Cfg != nil {
		item.CD = int64(item.Cfg.CD * 1000)
	}
	item.TimeCDFinish = t.provider.NowMs() + item.CD
}

func (t *Table) changeLevel(item *Item, level int) {
	item.Level = level
	// do something
}

func (t *Table) SetNormalAttackSkill(id common.SkillId) {
	item := allocItem()
	t.initItem(item, id, 1)
	item.TimeCDFinish = 0
	t.normalAttack = item
}

func (t *Table) findIndex(id common.SkillId) int {
	return slices.IndexFunc(t.skills, func(item *Item) bool {
		return item.SkillId == id
	})
}

func (t *Table) findEmptySlot() int {
	return t.findIndex("")
}

func (t *Table) getSkillItem(id common.SkillId) *Item {
	index := t.findIndex(id)
	if index == -1 {
		return nil
	}
	return t.skills[index]
}

func (t *Table) AddSkill(id common.SkillId, level int) {
	index := t.findIndex(id)
	if index != -1 {
		item := t.skills[index]
		t.changeLevel(item, level)
		return
	}

	item := &Item{}
	t.initItem(item, id, level)
	t.skills = append(t.skills, item)
}

func (t *Table) UpgradeSkill(id common.SkillId, level int) {
	index := t.findIndex(id)
	if index != -1 {
		item := t.skills[index]
		t.changeLevel(item, level)
	}
}

func (t *Table) RemoveSkill(id common.SkillId) {
	index := t.findIndex(id)
	if index == -1 {
		return
	}
	t.skills = slices.Delete(t.skills, index, index+1)
}

//func (t *Table) SetSkill(index int, id common.SkillId, level int) bool {
//	if !t.sureSize(index + 1) {
//		return false
//	}
//
//	indexOld := t.findIndex(id)
//	if indexOld != -1 {
//		//
//		return false
//	}
//
//	item := t.skills[index]
//	t.initItem(item, id, level)
//	return true
//}
//
//func (t *Table) sureSize(size int) bool {
//	if size <= 0 || size > MaxSlotNum {
//		return false
//	}
//
//	curSize := len(t.skills)
//	if curSize >= size {
//		return true
//	}
//
//	newData := make([]*Item, size-curSize)
//	for i := 0; i < size-curSize; i++ {
//		newData[i] = allocItem()
//	}
//	t.skills = append(t.skills, newData...)
//	return true
//}

func (t *Table) Update() {
	t.selectNextSkill()
}

func (t *Table) selectNextSkill() {
	t.nextSkill = nil

	// select cd ready skill
	t.nextSkill = t.selectFromSkills()
	if t.nextSkill != nil {
		return
	}

	// otherwise normal attack
	if t.IsNormalAttackReady() && !t.owner.HasSpecialStatus(common.SSNoNormalAttack) {
		t.nextSkill = t.normalAttack
	}
}

func (t *Table) selectFromSkills() *Item {
	if t.owner.HasSpecialStatus(common.SSNoSkill) {
		return nil
	}
	for _, v := range t.skills {
		if !t.IsItemCostReady(v) || !t.IsItemCDReady(v) {
			continue
		}
		return v
	}
	return nil
}

func (t *Table) GetNextSkill() (skillId common.SkillId, level int) {
	next := t.nextSkill
	if next == nil {
		return "", 1
	}
	return next.SkillId, next.Level
}

func (t *Table) PushCD(id common.SkillId, prefireTime int32) {
	if t.IsNormalAttack(id) {
		t.PushNormalAttackCD(prefireTime)
		return
	}
	t.pushCD(id)
}

func (t *Table) pushCD(id common.SkillId) {
	item := t.getSkillItem(id)
	if item == nil {
		return
	}
	item.TimeCDFinish = t.provider.NowMs() + item.CD
}

func (t *Table) getItemCDRest(item *Item) float32 {
	rest := item.TimeCDFinish - t.provider.NowMs()
	if rest < 0 {
		rest = 0
	}
	return float32(rest) / 1000
}

func (t *Table) GetCDRest(id common.SkillId) float32 {
	item := t.getSkillItem(id)
	if item == nil {
		return 0
	}
	return t.getItemCDRest(item)
}

func (t *Table) GetCDRestPercent(id common.SkillId) float32 {
	item := t.getSkillItem(id)
	if item == nil {
		return 1
	}

	if item.CD == 0 {
		return 0
	}
	return t.getItemCDRest(item) * 1000 / float32(item.CD)
}

func (t *Table) OffsetCD(id common.SkillId, offset float32) {
	item := t.getSkillItem(id)
	if item == nil {
		return
	}
	item.TimeCDFinish += int64(offset * 1000)
}

func (t *Table) OffsetCDPercent(id common.SkillId, offset float32) {
	item := t.getSkillItem(id)
	if item == nil {
		return
	}
	item.TimeCDFinish += int64(offset * float32(item.CD))
}

func (t *Table) IsCDReady(skillId common.SkillId) bool {
	item := t.getSkillItem(skillId)
	if item == nil {
		return false
	}
	return t.IsItemCDReady(item)
}

func (t *Table) IsItemCDReady(item *Item) bool {
	// cd is ready
	return t.provider.NowMs() >= item.TimeCDFinish
}

func (t *Table) IsItemCostReady(item *Item) bool {
	return t.owner.GetValue(common.Energy) >= 100
}

func (t *Table) SetNormalAttackInterval(interval float32) {
	old := t.normalAttackInterval
	t.normalAttackInterval = int64(interval * 1000)
	if old == 0 {
		return
	}

	// 速度变化了，更新一下normalattack的cd，应用攻速变化
	if t.normalAttack == nil {
		return
	}
	item := t.normalAttack
	rest := item.TimeCDFinish - t.provider.NowMs()
	// 已经无需更新
	if rest <= 0 {
		return
	}

	factor := float32(t.normalAttackInterval) / float32(old)
	final := float32(rest) * factor
	item.TimeCDFinish = t.provider.NowMs() + int64(final)
}

func (t *Table) IsNormalAttack(id common.SkillId) bool {
	item := t.normalAttack
	if item == nil {
		return false
	}
	return item.SkillId == id
}

func (t *Table) IsNormalAttackReady() bool {
	item := t.normalAttack
	if item == nil {
		return false
	}
	return t.provider.NowMs() >= item.TimeCDFinish
}

func (t *Table) PushNormalAttackCD(prefireTime int32) {
	item := t.normalAttack
	if item == nil {
		return
	}
	// 减去前摇时间
	now := t.provider.NowMs()
	item.TimeCDFinish = now + t.normalAttackInterval - int64(prefireTime)
}
