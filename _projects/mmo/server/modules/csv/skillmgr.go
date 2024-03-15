package csv

import (
	"github.com/dfklegend/cell2/utils/csvutils"

	"mmo/modules/csv/base"
	"mmo/modules/csv/entry"
)

// ----

type skillCfg struct {
	*base.DataMgr[*entry.Skill]
}

func NewSkillCfg() *skillCfg {
	return &skillCfg{
		DataMgr: base.NewDataMgr[*entry.Skill](),
	}
}

// ----

type skillEffectCfg struct {
	buf []*entry.SkillEffect
}

func NewSkillEffectCfg() *skillEffectCfg {
	return &skillEffectCfg{
		buf: []*entry.SkillEffect{},
	}
}

func (m *skillEffectCfg) Visit(visitor func(entry *entry.SkillEffect)) {
	for _, v := range m.buf {
		if v.SkillId != "" {
			visitor(v)
		}
	}
}

func (m *skillEffectCfg) LoadFromFile(path string) {
	err := csvutils.LoadFromFile(path, &m.buf)
	if err != nil {
		return
	}
}

// ----

type Effects struct {
	Effects []*entry.SkillEffect
}

func (e *Effects) Add(one *entry.SkillEffect) {
	e.Effects = append(e.Effects, one)
}

// skillEffects
// 技能: effects
type skillEffects struct {
	entries map[string]*Effects
}

func NewSkillEffects() *skillEffects {
	return &skillEffects{
		entries: map[string]*Effects{},
	}
}

// Build
// call after skillcfg skilleffectcfg loaded
func (m *skillEffects) Build(skill *skillCfg, effect *skillEffectCfg) {
	effect.Visit(func(entry *entry.SkillEffect) {
		m.addEffect(entry)
	})
}

func (m *skillEffects) addEffect(one *entry.SkillEffect) {
	effects := m.entries[one.SkillId]
	if effects == nil {
		effects = &Effects{
			Effects: []*entry.SkillEffect{},
		}
		m.entries[one.SkillId] = effects
	}

	effects.Add(one)
}

func (m *skillEffects) GetEntry(skillId string) *Effects {
	return m.entries[skillId]
}
