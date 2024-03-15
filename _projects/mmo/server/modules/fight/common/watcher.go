package common

import (
	"mmo/modules/fight/attr"
)

// IWatcher 战斗系统需要被同步出去
// 我们通过IWatcher来让外部场景模块来构建消息
// 场景自然的向视野内同步
type IWatcher interface {
	OnSkillStart(refChar ICharacter, data *DataSkillStart)
	OnSkillHit(refChar ICharacter, data *DataSkillHit)
	OnSkillBroken(refChar ICharacter, data *DataSkillBroken)
	OnAttrsChanged(refChar ICharacter, data *DataAttrsChanged)
}

// ---- 定义watcher的相关数据

type DataSkillStart struct {
	SkillId SkillId
	Src     CharId
	Tar     CharId
	Pos     Pos
}

type DataSkillHit struct {
	SkillId  SkillId
	Src      CharId
	Tar      CharId
	Dmg      int
	HPTar    int
	Critical bool
}

type DataToImpl struct {
}

type OneAttr struct {
	Index int
	Value attr.Value
}

type DataAttrsChanged struct {
	Attrs []*OneAttr
}

type DataSkillBroken struct {
	SkillId SkillId
	Src     CharId
}

// ----
