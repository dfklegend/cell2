package components

import (
	"math/rand"

	"github.com/dfklegend/cell2/utils/common"

	"mmo/common/entity"
	"mmo/messages/cproto"
	common2 "mmo/modules/fight/common"
	define2 "mmo/servers/scene/define"
	"mmo/servers/scene/entity/define"
)

const (
	SkillInit int32 = iota
	SkillPrefire
	SkillPostfire
)

// SkillComponent 负责技能执行
// 技能1s, 0.5s hit
type SkillComponent struct {
	*BaseSceneComponent

	self       *BaseUnit
	tarId      entity.EntityID
	tar        *BaseUnit
	skillId    int32
	timeStart  int64
	skillState int32
}

func NewSkillComponent() *SkillComponent {
	return &SkillComponent{
		BaseSceneComponent: NewBaseSceneComponent(),
	}
}

func (c *SkillComponent) OnPrepare() {
	c.BaseSceneComponent.OnPrepare()

	c.self = c.GetOwner().GetComponent(define.BaseUnit).(*BaseUnit)
}

func (c *SkillComponent) OnStart() {
}

func (c *SkillComponent) Update() {
	c.updateSkill()
}

func (c *SkillComponent) OnDestroy() {
}

func (c *SkillComponent) updateSkill() {
	if c.skillState == SkillInit {
		return
	}
	if c.skillState == SkillPrefire {
		now := common.NowMs()
		if now >= c.timeStart+500 {
			c.onSkillHit()
			c.skillState = SkillPostfire
		}
	}
	if c.skillState == SkillPostfire {
		now := common.NowMs()
		if now >= c.timeStart+1000 {
			c.skillState = SkillInit
		}
	}
}

func (c *SkillComponent) IsSkillRunning() bool {
	return c.skillState != SkillInit
}

func (c *SkillComponent) StartSkill(skillId int32, tarId entity.EntityID) {
	if c.IsSkillRunning() {
		return
	}
	c.skillId = skillId
	c.tarId = tarId
	c.timeStart = common.NowMs()
	c.skillState = SkillPrefire

	c.scene.PushViewMsg(define2.Pos{}, "startskill", &cproto.StartSklill{
		Id:      int32(c.GetOwner().GetId()),
		SkillId: "",
		Tar:     int32(c.tarId),
	})
}

func (c *SkillComponent) onSkillHit() {

	tarE := c.GetOwner().GetWorld().GetEntity(c.tarId)
	if tarE == nil {
		return
	}
	c.tar = tarE.GetComponent(define.BaseUnit).(*BaseUnit)
	if c.tar == nil {
		return
	}

	self := c.self
	tar := c.tar
	dmg := int32(self.GetChar().GetIntValue(common2.PhysicPower) - tar.GetChar().GetIntValue(common2.PhysicArmor))
	if dmg < 0 {
		dmg = 0
	}

	// 30% 概率高级
	critical := false
	if rand.Float32() < 0.3 {
		critical = true
		dmg *= 2
	}

	tar.ApplyDmg(dmg)

	c.scene.PushViewMsg(define2.Pos{}, "skillhit", &cproto.SkillHit{
		Id:       int32(c.GetOwner().GetId()),
		SkillId:  "",
		Tar:      int32(c.tarId),
		Dmg:      dmg,
		HPTar:    int32(tar.GetChar().GetHP()),
		Critical: critical,
	})
}
