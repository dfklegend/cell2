package test

import (
	"mmo/modules/fight/common"
)

type TestWatcher struct {
	now int64
}

func newTestWatcher() *TestWatcher {
	return &TestWatcher{}
}

func (p *TestWatcher) SetNow(now int64) {
}

func (p *TestWatcher) OnSkillStart(refChar common.ICharacter, data *common.DataSkillStart) {

}

func (p *TestWatcher) OnSkillHit(refChar common.ICharacter, data *common.DataToImpl) {

}

func (p *TestWatcher) OnSkillBreak(refChar common.ICharacter, data *common.DataToImpl) {

}
