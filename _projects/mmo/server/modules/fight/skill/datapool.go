package skill

import (
	"sync"

	"mmo/modules/fight/common"
)

var (
	dmgPool = sync.Pool{
		New: func() interface{} {
			return &common.DmgInstance{}
		},
	}

	skillStartPool = sync.Pool{
		New: func() interface{} {
			return &common.DataSkillStart{}
		},
	}

	skillHitPool = sync.Pool{
		New: func() interface{} {
			return &common.DataSkillHit{}
		},
	}
)

func AllocDmgInstance() *common.DmgInstance {
	//return &define.DmgInstance{}
	return dmgPool.Get().(*common.DmgInstance)
}

func FreeDmgInstance(data *common.DmgInstance) {
	data.Reset()
	dmgPool.Put(data)
}

func AllocDataSkillStart() *common.DataSkillStart {
	//return &common.DataSkillStart{}
	return skillStartPool.Get().(*common.DataSkillStart)
}

func FreeDataSkillStart(data *common.DataSkillStart) {
	skillStartPool.Put(data)
}

func AllocDataSkillHit() *common.DataSkillHit {
	//return &common.DataSkillHit{}
	return skillHitPool.Get().(*common.DataSkillHit)
}

func FreeDataSkillHit(data *common.DataSkillHit) {
	skillHitPool.Put(data)
}
