package buf

import (
	"mmo/modules/fight/common"
	"mmo/modules/fight/script"
)

const (
	rbScriptName = "鬼刀_机制"
	rbBufIdToAdd = "鬼刀_叠加"
)

func init() {
	mgr.Register(rbScriptName, func() script.IBufScript {
		return &rageBladeScript{}
	})
}

/* rageBladeScript
鬼索的狂暴之刃
每次普攻，攻速增加6%
+10% 攻速
+10% 法强
*/

type rageBladeScript struct {
	owner common.ICharProxy
}

func (s *rageBladeScript) Init(args ...any) {
}

func (s *rageBladeScript) OnStart(proxy script.IBufProxy) {
	owner := proxy.Owner()
	s.owner = owner
	owner.GetEvents().SubscribeWithReceiver("onSkillHit", s, s.onSkillHit)
}

func (s *rageBladeScript) OnEnd(proxy script.IBufProxy) {
	owner := proxy.Owner()
	owner.GetEvents().UnsubscribeWithReceiver("onSkillHit", s, s.onSkillHit)
}

func (s *rageBladeScript) OnTriggle(proxy script.IBufProxy) {
}

func (s *rageBladeScript) onSkillHit(args ...any) {
	if len(args) == 0 {
		return
	}
	skill := args[0].(script.ISkillProxy)
	if skill.IsBGSkill() {
		return
	}
	s.owner.AddBuf(rbBufIdToAdd, 1, 1)
}
