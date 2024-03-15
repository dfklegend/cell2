package skill

import (
	"mmo/common/config"
	"mmo/modules/csv"
	"mmo/modules/csv/entry"
	"mmo/modules/fight/common"
	"mmo/modules/fight/common/skilldef"
	skilleffect2 "mmo/modules/fight/common/skilleffect"
	"mmo/modules/fight/lua/env"
	"mmo/modules/fight/skill/define"
	factory2 "mmo/modules/fight/skill/effect/factory"
	"mmo/modules/fight/skill/formula/factory"
)

/*
	prefire->hit->postfire
	start
		Cancel		自行取消
		Break		被打断

		onEnd		正常结束
		onFailed	失败

	现在技能在释放时，扣cd和消耗
*/

type Skill struct {
	id common.SkillId
	// 某些技能可能改变执行行为，变成另外一个技能
	// 此时, cd还是原有技能
	idShell common.SkillId
	owner   common.ICharacter

	timeProvider common.ITimeProvider
	watcher      common.IWatcher
	world        common.IWorld

	script *Script
	proxy  *Proxy

	cfg       *entry.Skill
	effects   *csv.Effects
	level     int
	startTime int64
	state     int

	hitTime   int32
	totalTime int32

	// 目标
	tar    common.CharId
	tarPos common.Pos

	// 是后台技能(触发技能)
	bgSkill bool
	// 技能释放失败，比如被打断
	failedReason int

	subSkill        bool
	subRunning      *Skill
	subRunningIndex int
	stackDepth      int
}

func NewSkill(id common.SkillId, level int, cfg *entry.Skill,
	owner common.ICharacter) *Skill {
	s := &Skill{
		id:       id,
		level:    level,
		cfg:      cfg,
		owner:    owner,
		bgSkill:  false,
		subSkill: false,
	}

	if owner != nil {
		s.SetWorld(owner.GetWorld())
	}
	s.effects = csv.SkillEffects.GetEntry(id)
	s.proxy = newProxy(s)
	return s
}

// prepare

func (s *Skill) SetTar(id common.CharId) {
	s.tar = id
}

func (s *Skill) SetWorld(world common.IWorld) {
	s.world = world
	if world != nil {
		s.timeProvider = world.GetTimeProvider()
		s.watcher = world.GetWatcher()
	}
}

// prepare over

func (s *Skill) SetBG(bg bool) {
	s.bgSkill = bg
}

func (s *Skill) SetSubSkill(v bool) {
	s.subSkill = v
}

func (s *Skill) IsBGSkill() bool {
	return s.bgSkill
}

func (s *Skill) GetCfg() *entry.Skill {
	return s.cfg
}

func (s *Skill) setState(state int) {
	s.state = state
}

func (s *Skill) IsOver() bool {
	return s.state == Over || s.state == Failed
}

func (s *Skill) IsFailed() bool {
	return s.state == Failed
}

func (s *Skill) GetFailedReason() int {
	return s.failedReason
}

// Start
// 开始技能
func (s *Skill) Start() {
	s.startTime = s.timeProvider.NowMs()
	s.setState(Prefire)
	s.initByTime(0)

	s.tryFindScript()

	if s.watcher != nil {
		data := AllocDataSkillStart()
		data.SkillId = s.id
		data.Src = s.owner.GetId()
		data.Tar = s.tar

		s.watcher.OnSkillStart(s.owner, data)
		FreeDataSkillStart(data)
	}

	recorder := s.owner.GetWorld().GetDetailRecorder()
	if recorder != nil {
		tar := s.world.GetChar(s.tar)
		recorder.OnStartSkill(s.id, s.owner, tar)
	}

	s.scriptDoStart()

	s.pushCD()
	s.applyCost()

	// 比如: 开始释放技能就套个buf
	s.applyEffects(skilleffect2.CastCheck, nil)

	if s.hitTime == 0 {
		s.doHit()
	}
}

func (s *Skill) tryFindScript() {
	if config.EnableFightLuaPlugin == 0 || s.world.GetLua() == nil {
		return
	}

	env := s.world.GetLua().GetEnvData().(*env.ScriptEnvData)
	mgr := env.GetValue("skillScripts").(*ScriptMgr)
	s.script = mgr.GetScript(s.id)
}

// 普攻加速规则
// 如果攻击间隔大于动作时间，保持不变
// 如果攻击间隔小于动作时间，动作加速
func (s *Skill) initTime() {
	var tarTime int32
	if s.cfg.NormalAttack {
		tarTime = s.owner.GetNormalAttackInterval()
	}
	s.initByTime(tarTime)
}

func (s *Skill) initByTime(tarTime int32) {
	s.hitTime = s.cfg.HitTime
	s.totalTime = s.cfg.TotalTime

	if tarTime <= 0 || s.totalTime < tarTime {
		return
	}

	if s.totalTime == 0 {
		s.totalTime = 1
	}

	// scale
	scale := float32(tarTime) / float32(s.totalTime)
	s.hitTime = int32(scale * float32(s.hitTime))
	s.totalTime = tarTime
}

func (s *Skill) Update() {
	// check hit
	switch s.state {
	case Prefire:
		s.updatePrefire()
	case Postfire:
		s.updatePostfire()
	case SubSkillRunning:
		s.updateSubSkills()
	}
}

// Cancel
// 自行取消
func (s *Skill) Cancel() {
	switch s.state {
	case Init:
	case Over:
		return
	}
	s.switchToFailed(ReasonCancel)
}

// Break
// 被打断，比如眩晕等
func (s *Skill) Break() {
	switch s.state {
	case Init:
	case Over:
		return
	}

	s.switchToFailed(ReasonBreak)
}

func (s *Skill) switchToFailed(reason int) {
	if s.state == SubSkillRunning && s.subRunning != nil {
		// break child
		s.subRunning.Break()
	}

	s.setState(Failed)
	s.failedReason = reason

	s.onFailed()
	s.clear()
}

func (s *Skill) onFailed() {
	s.applyEffects(skilleffect2.OnFailedCheck, nil)
}

func (s *Skill) updatePrefire() {
	now := s.timeProvider.NowMs()
	if now < s.startTime+int64(s.hitTime) {
		return
	}
	s.doHit()
	s.setState(Postfire)
}

func (s *Skill) updatePostfire() {
	now := s.timeProvider.NowMs()
	if now < s.startTime+int64(s.totalTime) {
		return
	}
	s.doEnd()
}

func (s *Skill) doEnd() {
	s.onEnd()
	s.setState(Over)
}

func (s *Skill) getFormula() define.IFormula {
	obj := factory.GetFormulaFactory().Create(skilldef.FromulaNormal)
	if obj == nil {
		return nil
	}
	return obj.(define.IFormula)
}

func (s *Skill) doHit() {
	s.applyHit()
	// 如果有组合技能，执行子技能
	if s.hasSubSkills() && s.checkCanBeginSubSkills() {
		s.setState(SubSkillRunning)
		s.beginSubSkills()
	} else {
		if s.totalTime > s.hitTime {
			s.setState(Postfire)
		} else {
			// 没有后摇，直接结束了
			s.doEnd()
		}
	}
}

func (s *Skill) applyHit() {
	if !s.HasValidTar() {
		s.Cancel()
		return
	}

	s.pushCD()

	// todo: 未来，根据技能类型来执行
	formula := s.getFormula()
	if formula == nil {
		return
	}
	dmg := AllocDmgInstance()

	tar := s.world.GetChar(s.tar)
	// 结算
	// 如果要结算命中，命中通过之后才
	s.applyEffects(skilleffect2.HitCasterCheck, tar)
	if tar != nil {
		s.applyEffects(skilleffect2.HitTargetCheck, tar)
		formula.Apply(s, s.owner, tar, dmg)
	}

	//
	if s.watcher != nil {
		data := AllocDataSkillHit()
		data.SkillId = s.id
		data.Src = s.owner.GetId()
		data.Tar = s.tar
		data.Dmg = dmg.Dmg
		data.Critical = dmg.Critical
		if tar != nil {
			data.HPTar = tar.GetHP()
		}

		s.watcher.OnSkillHit(s.owner, data)
		FreeDataSkillHit(data)
	}

	recorder := s.owner.GetWorld().GetDetailRecorder()
	if recorder != nil {
		recorder.OnSkillHit(s.id, s.owner, tar, dmg.Dmg)
	}

	s.scriptDoHit()
	s.owner.OnSkillHit(dmg)
	s.scriptDoBeHit(tar)

	if tar != nil {
		tar.OnSkillBeHit(dmg)

		if dmg.TarBeKilled {
			s.owner.OnKillTarget(tar.GetId(), s.id)
		}
	}

	FreeDmgInstance(dmg)
}

func (s *Skill) pushCD() {
	id := s.id
	if s.idShell != "" {
		id = s.idShell
	}
	s.owner.PushSkillCD(id, 0)
}

func (s *Skill) applyCost() {
	if s.IsBGSkill() || s.cfg.NormalAttack || s.subSkill {
		return
	}
	s.owner.ApplyMPCost(100)
}

func (s *Skill) scriptDoStart() {
	if s.script != nil {
		s.script.onSkillStart(s.proxy)
	}
}

func (s *Skill) scriptDoHit() {
	if s.script != nil {
		s.script.onSkillHit(s.proxy)
	}
	s.owner.GetEvents().Publish("onSkillHit", s.proxy)
}

func (s *Skill) scriptDoBeHit(tar common.ICharacter) {
	if tar == nil {
		return
	}
	tar.GetEvents().Publish("onSkillBeHit", s.proxy)
}

// collectTargets
// 搜集实际的目标，比如，范围技能
// 然后对每个目标执行结算
func (s *Skill) collectTargets() {
}

// 正常结束
func (s *Skill) onEnd() {
	s.clear()
}

// applyEffects 在技能不同阶段，效用技能效果
func (s *Skill) applyEffects(time int, tar common.ICharacter) {
	if s.effects == nil {
		return
	}

	// for debug
	//if s.id == "盖伦技能" && time == skilleffect2.CastCheck {
	//	l.L.Info("hh")
	//}

	for _, v := range s.effects.Effects {
		if v.ApplyTime.Type != time {
			continue
		}
		s.applyEffect(tar, v)
	}
}

func (s *Skill) getEffect(op int) define.ISkillEffect {
	obj := factory2.GetOpFactory().Create(op)
	if obj == nil {
		return nil
	}
	return obj.(define.ISkillEffect)
}

func (s *Skill) applyEffect(tar common.ICharacter, effect *entry.SkillEffect) {
	op := s.getEffect(effect.Op.Type)
	if op == nil {
		return
	}

	finalTar := tar
	if effect.TarType.Type == skilleffect2.TarCaster {
		finalTar = s.owner
	}

	op.Apply(s.owner, finalTar, effect, 1)
}

// clear做对象清理
func (s *Skill) clear() {

}

func (s *Skill) HasValidTar() bool {
	return true
}
