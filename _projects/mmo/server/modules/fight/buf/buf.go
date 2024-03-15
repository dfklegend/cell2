package buf

import (
	"mmo/modules/fight/common"
	"mmo/modules/fight/common/bufdef"
	"mmo/modules/fight/script"

	"mmo/modules/csv/entry"
	"mmo/modules/csv/structs"
)

const (
	neverOver = -1
)

type IBuf interface {
	Start(caster common.ICharacter, owner common.ICharacter, cfg *entry.Buf, level, stack int)

	GetId() common.BufId
	GetLevel() int
	GetStack() int

	Refresh(level, stack int)

	Cancel()
	Update()
	IsOver() bool
}

type Buf struct {
	owner    common.ICharacter
	casterId common.CharId
	id       common.BufId
	level    int

	cfg   *entry.Buf
	times int

	timeProvider common.ITimeProvider

	triggled        int
	nextTriggleTime int64
	totalTime       int64

	stack    int
	canceled bool

	proxy  *Proxy
	script script.IBufScript
}

func NewBuf() *Buf {
	return &Buf{
		canceled:  false,
		triggled:  0,
		totalTime: -1,
	}
}

func (b *Buf) GetId() common.BufId {
	return b.id
}

func (b *Buf) GetLevel() int {
	return b.level
}

func (b *Buf) GetStack() int {
	return b.stack
}

func (b *Buf) Start(caster common.ICharacter, owner common.ICharacter, cfg *entry.Buf, level, stack int) {
	if caster == nil || owner == nil || cfg == nil || level < 1 || stack < 1 {
		return
	}
	b.owner = owner
	b.casterId = caster.GetId()
	b.id = cfg.Id
	b.level = level
	b.cfg = cfg
	b.times = cfg.Times
	b.stack = b.clampStack(stack)

	b.proxy = newProxy(b)
	b.tryCreateScript()

	b.timeProvider = caster.GetWorld().GetTimeProvider()
	b.initTime()

	if !b.isTriggledOver() {
		b.doStart()
	}
}

func (b *Buf) doStart() {
	b.offsetAttrs(b.owner, true, b.level, b.stack)
	b.offsetSpecialStatus(b.owner, true)
	if b.script != nil {
		b.script.OnStart(b.proxy)
	}
}

func (b *Buf) doEnd() {
	b.offsetAttrs(b.owner, false, b.level, b.stack)
	b.offsetSpecialStatus(b.owner, false)

	if b.script != nil {
		b.script.OnEnd(b.proxy)
		b.script = nil
	}
}

// cfg.Times == -1 代表永久存在
func (b *Buf) initTime() {
	cfg := b.cfg

	if cfg.Interval > 0 {
		provider := b.timeProvider
		b.nextTriggleTime = provider.NowMs() + int64(cfg.Interval)
	}

	if cfg.Times > 0 {
		b.totalTime = int64(cfg.Times * cfg.Interval)
	}
}

func (b *Buf) IsOver() bool {
	return b.isTriggledOver() || b.canceled
}

func (b *Buf) isTriggledOver() bool {
	if b.times == neverOver {
		return false
	}
	return b.triggled >= b.times
}

func (b *Buf) offsetAttrs(tar common.ICharacter, add bool, level, stack int) {
	cfg := b.cfg
	if cfg == nil {
		return
	}

	if cfg.Attr0.Type != -1 {
		b.offsetAttr(tar, add, level, stack, cfg.Attr0, cfg.Base0, cfg.Lv0)
	}
	if cfg.Attr1.Type != -1 {
		b.offsetAttr(tar, add, level, stack, cfg.Attr1, cfg.Base1, cfg.Lv1)
	}
	if cfg.Attr2.Type != -1 {
		b.offsetAttr(tar, add, level, stack, cfg.Attr2, cfg.Base2, cfg.Lv2)
	}
}

func (b *Buf) offsetAttr(tar common.ICharacter, add bool, level, stack int, attr structs.AttrValue, base, lv float32) {
	if attr.Type == -1 {
		return
	}

	final := (base + lv*float32(level-1)) * float32(stack)

	if !add {
		final = -final
	}

	if attr.IsPercent {
		tar.OffsetPercent(attr.Type, final)
		return
	}
	tar.OffsetBase(attr.Type, float64(final))
}

func (b *Buf) offsetSpecialStatus(tar common.ICharacter, add bool) {
	cfg := b.cfg
	if cfg == nil || cfg.SpecialStatus.Ints == nil {
		return
	}

	for i := 0; i < len(cfg.SpecialStatus.Ints); i++ {
		one := cfg.SpecialStatus.Ints[i]
		if one < 0 {
			continue
		}
		if add {
			tar.AddSpecialStatus(one)
		} else {
			tar.SubSpecialStatus(one)
		}
	}
}

func (b *Buf) Update() {
	if b.cfg.Interval > 0 {
		b.updateTriggle()
	}
}

func (b *Buf) updateTriggle() {
	now := b.timeProvider.NowMs()
	if now >= b.nextTriggleTime {
		b.triggle()
		b.nextTriggleTime = now + int64(b.cfg.Interval)
	}

	if b.isTriggledOver() {
		b.doEnd()
	}
}

func (b *Buf) triggle() {
	b.triggled++
	cfg := b.cfg
	if cfg.TriggleSkillId != "" {
		b.doTriggleSkill()
	}
	// call triggle
	if b.script != nil {
		b.script.OnTriggle(b.proxy)
	}
}

// 比如: 给owner加血, 伤害, 释放aoe等
func (b *Buf) selectSkillTar() common.ICharacter {
	cfg := b.cfg
	if cfg.SkillTar.Type == bufdef.TarTypeOwnerTar {
		return b.owner.GetTar()
	}
	return b.owner
}

func (b *Buf) doTriggleSkill() {
	cfg := b.cfg
	tar := b.selectSkillTar()
	if tar != nil {
		b.owner.CallbackSkill(cfg.TriggleSkillId, 1, b.owner, tar.GetId())
	}
}

func (b *Buf) Cancel() {
	if b.IsOver() {
		return
	}
	b.canceled = true
	b.doEnd()
}

func (b *Buf) clampStack(stack int) int {
	if stack <= int(b.cfg.MaxStack) {
		return stack
	}
	return int(b.cfg.MaxStack)
}

func (b *Buf) offsetCurAttrs(add bool) {
	b.offsetAttrs(b.owner, add, b.level, b.stack)
}

func (b *Buf) Refresh(level, stack int) {
	if level < b.level {
		// 低级buf, 抛弃掉
		return
	}
	// buf 刷新
	if level > b.level {
		// 高等级buf
		b.offsetCurAttrs(false)

		b.level = level
		b.stack = b.clampStack(stack)
		b.offsetCurAttrs(true)
		b.refreshTriggle()
		return
	}

	// 等级相同，叠加层数，同时刷新
	b.offsetCurAttrs(false)
	b.stack = b.clampStack(b.stack + stack)
	b.offsetCurAttrs(true)
	b.refreshTriggle()
}

// TODO: 是否需要明确设置一个刷新类型
// 如果必刷新nextTriggleTime, 可能导致，比如周期性回血，一直被刷新，触发不了回血
func (b *Buf) refreshTriggle() {
	b.triggled = 0
	freqTriggle := b.cfg.Times > 1
	if !freqTriggle {
		// 时间也刷一下
		b.refreshNextTriggleTime()
	}
}

func (b *Buf) refreshNextTriggleTime() {
	b.nextTriggleTime = b.timeProvider.NowMs() + int64(b.cfg.Interval)
}
