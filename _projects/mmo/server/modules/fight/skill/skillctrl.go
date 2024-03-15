package skill

import (
	"golang.org/x/exp/slices"

	"mmo/modules/csv"
	"mmo/modules/fight/common"
	"mmo/modules/fight/lua/env"
)

// Ctrl
// 控制技能
// 当前有一个主技能
// 可以有很多后台技能
type Ctrl struct {
	main *Skill
	bg   []*Skill
}

func NewCtrl() *Ctrl {
	return &Ctrl{
		bg: []*Skill{},
	}
}

func (c *Ctrl) Update() {
	if c.main != nil {
		if c.updateSkill(c.main) {
			c.main = nil
		}
	}
	c.updateBG()
}

func (c *Ctrl) updateBG() {
	for i := 0; i < len(c.bg); {
		skill := c.bg[i]
		if c.updateSkill(skill) {
			c.bg = slices.Delete(c.bg, i, i+1)
		} else {
			i++
		}
	}
}

func (c *Ctrl) IsSkillRunning() bool {
	return c.main != nil
}

//	return true, skill is over
func (c *Ctrl) updateSkill(skill *Skill) bool {
	if !skill.IsOver() {
		skill.Update()
		return false
	}
	return true
}

func (c *Ctrl) CastSkill(id common.SkillId, level int, src common.ICharacter, tarId common.CharId) {
	skill := c.checkNewSkill(id, level, src, tarId)
	if skill != nil {
		c.main = skill
		skill.Start()
	}
}

func (c *Ctrl) scriptChangeSkill(id common.SkillId, src common.ICharacter) common.SkillId {
	env := src.GetWorld().GetLua().GetEnvData().(*env.ScriptEnvData)
	mgr := env.GetValue("skillScripts").(*ScriptMgr)
	script := mgr.GetScript(id)
	if script == nil {
		return ""
	}
	return script.changeSkill(src)
}

func (c *Ctrl) newSkillToTarget(id common.SkillId, level int, src common.ICharacter, tarId common.CharId) *Skill {
	cfg := csv.Skill.GetEntry(id)
	if cfg == nil {
		return nil
	}
	skill := NewSkill(id, level, cfg, src)
	skill.SetTar(tarId)
	skill.SetWorld(src.GetWorld())
	return skill
}

func (c *Ctrl) CallbackSkill(id common.SkillId, level int, src common.ICharacter, tarId common.CharId) {
	skill := c.checkNewSkill(id, level, src, tarId)
	if skill != nil {
		skill.SetBG(true)
		skill.Start()

		if !skill.IsOver() {
			c.bg = append(c.bg, skill)
		}
	}
}

// checkNewSkill 检查是否需要改变实际执行技能
func (c *Ctrl) checkNewSkill(id common.SkillId, level int, src common.ICharacter, tarId common.CharId) *Skill {
	// 技能切换判定
	newId := c.scriptChangeSkill(id, src)
	changed := false
	if newId != "" {
		changed = true
	} else {
		newId = id
	}
	skill := c.newSkillToTarget(newId, level, src, tarId)
	if skill == nil {
		return nil
	}
	if changed {
		skill.idShell = id
	}
	return skill
}

func (c *Ctrl) BreakSkill(src common.ICharacter, breakNormalAttack bool) {
	// 主技能
	if c.main != nil {
		if c.tryBreak(c.main, breakNormalAttack) {
			src.OnSkillBroken(c.main.id)
			c.main = nil
		}
	}

	// 后台技能(现在一般瞬发的，先不处理)
}

func (c *Ctrl) tryBreak(s *Skill, breakNormalAttack bool) bool {
	if s.GetCfg().NormalAttack {
		if breakNormalAttack {
			s.Break()
			return true
		}
	}

	if breakNormalAttack {
		return false
	}

	s.Break()
	return true
}
