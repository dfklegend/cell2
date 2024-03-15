package sceneservice

import (
	"mmo/messages/cproto"
	common2 "mmo/modules/fight/common"
	define2 "mmo/servers/scene/define"
)

//

type FightWatcher struct {
	scene define2.IScene
}

func newFightWatcher(scene define2.IScene) *FightWatcher {
	return &FightWatcher{
		scene: scene,
	}
}

func (w *FightWatcher) OnSkillStart(refChar common2.ICharacter, data *common2.DataSkillStart) {
	w.scene.PushViewMsg(define2.Pos{}, "startskill", &cproto.StartSklill{
		Id:      int32(data.Src),
		SkillId: data.SkillId,
		Tar:     int32(data.Tar),
	})
}

func (w *FightWatcher) OnSkillHit(refChar common2.ICharacter, data *common2.DataSkillHit) {
	w.scene.PushViewMsg(define2.Pos{}, "skillhit", &cproto.SkillHit{
		Id:       int32(data.Src),
		SkillId:  data.SkillId,
		Tar:      int32(data.Tar),
		Dmg:      int32(data.Dmg),
		HPTar:    int32(data.HPTar),
		Critical: data.Critical,
	})
}

func (w *FightWatcher) OnSkillBroken(refChar common2.ICharacter, data *common2.DataSkillBroken) {
	w.scene.PushViewMsg(define2.Pos{}, "skillbroken", &cproto.SkillBroken{
		Id:      int32(data.Src),
		SkillId: data.SkillId,
	})
}

func (w *FightWatcher) OnAttrsChanged(refChar common2.ICharacter, data *common2.DataAttrsChanged) {
	attrs := make([]*cproto.CardAttr, 0)
	for _, v := range data.Attrs {
		attrs = append(attrs, &cproto.CardAttr{
			Index: int32(v.Index),
			Value: float32(v.Value),
		})
	}
	w.scene.PushViewMsg(define2.Pos{}, "attrschanged", &cproto.UnitAttrsChanged{
		Id:    int32(refChar.GetId()),
		Attrs: attrs,
	})
}
