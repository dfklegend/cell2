package effect

import (
	"math/rand"

	"mmo/modules/csv/base"
	"mmo/modules/csv/entry"
	"mmo/modules/fight/common"
	"mmo/modules/fight/common/skilleffect"
)

func init() {
	Register(skilleffect.OpAddBuf, &addBuf{})
}

type addBuf struct {
}

type addBufArgs struct {
	Rate  float32
	BufId string
	Stack int
}

func (a *addBuf) Format(cfg *base.IArgs) {
	cfg.FormatArgs(&addBufArgs{
		Rate:  0,
		BufId: "",
		Stack: 0,
	})
}

// rate, bufid, stack

func (a *addBuf) Apply(caster common.ICharacter, tar common.ICharacter, cfg *entry.SkillEffect, skillLv int) {
	if tar == nil {
		return
	}

	bufArgs := cfg.ArgsImpl.(*addBufArgs)
	rate := bufArgs.Rate
	bufId := bufArgs.BufId
	stack := bufArgs.Stack

	if bufId == "" || rate <= 0 || stack <= 0 {
		return
	}

	if rand.Float32() > rate {
		return
	}

	tar.AddBuf(caster, tar, bufId, skillLv, stack)
	//l.Log.Infof("%v add buf %v", tar.GetId(), bufId)
}
