package effect

import (
	"mmo/modules/csv/base"
	"mmo/modules/csv/entry"
	"mmo/modules/fight/common"
	"mmo/modules/fight/common/skilleffect"
)

func init() {
	Register(skilleffect.OpRemoveBuf, &removeBuf{})
}

type removeBuf struct {
}

type removeBufArgs struct {
	BufId string
}

func (a *removeBuf) Format(cfg *base.IArgs) {
	cfg.FormatArgs(&removeBufArgs{
		BufId: "",
	})
}

// bufid

func (a *removeBuf) Apply(caster common.ICharacter, tar common.ICharacter, cfg *entry.SkillEffect, skillLv int) {
	if tar == nil {
		return
	}

	bufArgs := cfg.ArgsImpl.(*removeBufArgs)
	bufId := bufArgs.BufId
	tar.RemoveBuf(bufId)
}
