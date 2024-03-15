package charimpls

import (
	"mmo/modules/fight/common"
)

func (c *Character) ApplyDmg(dmg int) int {
	if dmg <= 0 {
		return 0
	}
	if c.IsDead() {
		return 0
	}

	if c.IsInvincible() {
		return 0
	}

	hp := c.GetHP()
	oldHp := hp
	hp -= dmg

	if hp < 0 {
		hp = 0
	}
	c.SetHP(hp)
	if hp == 0 {
		c.onDead()
	}
	return oldHp - hp
}

func (c *Character) onDead() {
	if c.charWatcher != nil {
		c.charWatcher.OnDead()
	}
}

func (c *Character) IsDead() bool {
	return c.GetHP() == 0
}

func (c *Character) Born() {
	c.SetIntBaseValue(common.HP, c.GetIntValue(common.HPMax))
}

func (c *Character) IsSkillRunning() bool {
	return c.skill.IsSkillRunning()
}

func (c *Character) GetTarId() common.CharId {
	return c.idTar
}

func (c *Character) GetTar() common.ICharacter {
	return c.world.GetChar(c.idTar)
}

func (c *Character) CastSkill(id common.SkillId, level int, src common.ICharacter, tar common.CharId) {
	c.skill.CastSkill(id, level, src, tar)
}

func (c *Character) CallbackSkill(id common.SkillId, level int, src common.ICharacter, tar common.CharId) {
	c.skill.CallbackSkill(id, level, src, tar)
}

func (c *Character) BreakSkill(src common.ICharacter, breakNormalAttack bool) {
	c.skill.BreakSkill(src, breakNormalAttack)
}

func (c *Character) UpdateAttack(tar common.CharId) {
	c.idTar = tar
	// 如果当前技能
	if c.skill.IsSkillRunning() {
		return
	}
	c.skillTable.Update()
	skillId, level := c.skillTable.GetNextSkill()
	if skillId == "" {
		return
	}
	c.skill.CastSkill(skillId, level, c, tar)
}

func (c *Character) ClearTar() {
	c.idTar = -1
}

func (c *Character) PushSkillCD(id common.SkillId, prefireTime int32) {
	c.skillTable.PushCD(id, prefireTime)
}

func (c *Character) OnSkillHit(dmg *common.DmgInstance) {
	c.addEnergyAndBroadcast(5)
}

func (c *Character) OnSkillBeHit(dmg *common.DmgInstance) {
	c.addEnergyAndBroadcast(5)
}

func (c *Character) addEnergyAndBroadcast(off int) {
	c.AddEnergy(off)
	c.broadcastAttrValue(common.Energy)
}

func (c *Character) OnSkillBroken(id common.SkillId) {
	if c.world.GetWatcher() == nil {
		return
	}
	data := &common.DataSkillBroken{
		SkillId: id,
		Src:     c.GetId(),
	}

	c.world.GetWatcher().OnSkillBroken(c, data)
}

func (c *Character) OnKillTarget(tar common.CharId, skillId common.SkillId) {
	if c.charWatcher != nil {
		c.charWatcher.OnKillTarget(tar, skillId)
	}
}
