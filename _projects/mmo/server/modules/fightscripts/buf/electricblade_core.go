package buf

import (
	"mmo/modules/fight/common"
	"mmo/modules/fight/script"
)

const (
	ebScriptName = "电刀_机制"
	ebHitSkill   = "电刀伤害"
)

func init() {
	mgr.Register(ebScriptName, func() script.IBufScript {
		return &electricbladeScript{}
	})
}

// electricbladeScript
// 电刀
// 3次攻击一次电击伤害
type electricbladeScript struct {
	idSkillHit int
	tick       int
	owner      common.ICharProxy
}

func (s *electricbladeScript) Init(args ...any) {
}

func (s *electricbladeScript) OnStart(proxy script.IBufProxy) {
	owner := proxy.Owner()
	s.owner = owner
	owner.GetEvents().SubscribeWithReceiver("onSkillHit", s, s.onSkillHit)
}

func (s *electricbladeScript) OnEnd(proxy script.IBufProxy) {
	owner := proxy.Owner()
	owner.GetEvents().UnsubscribeWithReceiver("onSkillHit", s, s.onSkillHit)
}

func (s *electricbladeScript) OnTriggle(proxy script.IBufProxy) {
}

func (s *electricbladeScript) onSkillHit(args ...any) {
	if len(args) == 0 {
		return
	}
	skill := args[0].(script.ISkillProxy)
	if skill.IsBGSkill() {
		return
	}
	s.tick++
	if s.tick >= 3 {
		s.tick = 0
		owner := s.owner
		owner.GoCallbackSkill(ebHitSkill, 1, owner, owner.GetTarId())
	}
}
