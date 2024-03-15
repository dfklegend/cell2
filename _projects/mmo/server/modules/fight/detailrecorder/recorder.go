package detailrecorder

import (
	"fmt"

	"mmo/common/fightutils"
	"mmo/modules/fight/common"
)

type Recorder struct {
	infos        []string
	timeProvider common.ITimeProvider

	attrCWatcher *fightutils.AttrChangeWatcher
}

func NewRecorder(timeProvider common.ITimeProvider) *Recorder {
	return &Recorder{
		infos:        make([]string, 0, 1000),
		timeProvider: timeProvider,
		attrCWatcher: fightutils.NewAttrWatcher(),
	}
}

func (r *Recorder) addInfo(info string) {
	r.infos = append(r.infos, info)
}

func (r *Recorder) log(format string, args ...any) {
	args = append([]interface{}{r.timeProvider.NowMs()}, args...)
	r.addInfo(fmt.Sprintf("%v "+format, args...))
}

func (r *Recorder) OnStartFight() {
	r.log("战斗开始")
}

func (r *Recorder) OnEndFight() {
	r.log("战斗结束")
}

func (r *Recorder) Visit(doFunc func(info string)) {
	for _, v := range r.infos {
		doFunc(v)
	}
}

func (r *Recorder) OnStartSkill(id common.SkillId, src common.ICharacter, tar common.ICharacter) {
	r.log(" %v 对 %v 施放技能: %v", src.GetId(), tar.GetId(), id)
}

func (r *Recorder) OnSkillHit(id common.SkillId, src common.ICharacter, tar common.ICharacter, dmg int) {
	r.log(" %v 对 %v 命中技能: %v 伤害: %v 剩余血量: %v", src.GetId(), tar.GetId(), id, dmg, tar.GetHP())
}

func (r *Recorder) OnPreAddBuf(c common.ICharacter, id common.BufId, level, stack int) {
	r.attrCWatcher.Reset()
	c.AddAttrChangeWatcher(r.attrCWatcher)
}

func (r *Recorder) OnPostAddBuf(c common.ICharacter, id common.BufId, level, stack int) {
	r.log(" %v 获得buf: %v(lv:%v) stack: %v", c.GetId(), id, level, stack)
	r.attrCWatcher.Visit(func(one *fightutils.Attr) {
		r.log("  %v  %v -> %v", common.AttrIndexToName(one.Index), one.OldV, one.NewV)
	})
	c.RemoveAttrChangeWatcher(r.attrCWatcher)
}
