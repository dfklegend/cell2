package buf

import (
	l "github.com/dfklegend/cell2/utils/logger"
	"golang.org/x/exp/slices"

	"mmo/modules/csv"
	"mmo/modules/fight/common"
)

// Ctrl
// 添加和移除buf，已处理了同帧多次添加移除同一个buf的问题
type Ctrl struct {
	bufs []IBuf
}

func NewCtrl() *Ctrl {
	return &Ctrl{
		bufs: []IBuf{},
	}
}

func (c *Ctrl) AddBuf(caster common.ICharacter, owner common.ICharacter, id common.BufId, level int, stack int) {
	cfg := csv.Buf.GetEntry(id)
	if cfg == nil {
		l.L.Warnf("can not find buf: %v", id)
		return
	}

	recoder := owner.GetWorld().GetDetailRecorder()

	// 无敌免疫有害buf
	if owner.IsInvincible() && cfg.AffectType == common.BufAffectBad {
		return
	}

	// preaddbuf
	if recoder != nil {
		recoder.OnPreAddBuf(owner, id, level, stack)
	}
	one := c.doAddBuf(caster, owner, id, level, stack)
	// post add buf
	if recoder != nil {
		if one != nil {
			recoder.OnPostAddBuf(owner, id, level, one.GetStack())
		}
	}
}

func (c *Ctrl) doAddBuf(caster common.ICharacter, owner common.ICharacter, id common.BufId, level int, stack int) IBuf {
	cfg := csv.Buf.GetEntry(id)
	if cfg == nil {
		l.L.Warnf("can not find buf: %v", id)
		return nil
	}

	// 如果存在buf
	_, old := c.GetActiveBufObj(id)
	if old != nil {
		// 刷一下buf
		if level >= old.GetLevel() {
			old.Refresh(level, stack)
		}
		return old
	}

	buf := NewBuf()
	buf.Start(caster, owner, cfg, level, stack)
	c.bufs = append(c.bufs, buf)
	return buf
}

func (c *Ctrl) findActiveIndex(id common.BufId) int {
	return slices.IndexFunc(c.bufs, func(b IBuf) bool {
		return id == b.GetId() && !b.IsOver()
	})
}

func (c *Ctrl) RemoveBuf(id common.BufId) {
	_, buf := c.GetActiveBufObj(id)
	if buf == nil {
		return
	}
	buf.Cancel()
}

func (c *Ctrl) HasBuf(id common.BufId) bool {
	return c.findActiveIndex(id) != -1
}

func (c *Ctrl) GetActiveBufObj(id common.BufId) (int, IBuf) {
	index := c.findActiveIndex(id)
	if index == -1 {
		return -1, nil
	}
	return index, c.bufs[index]
}

func (c *Ctrl) GetBuf(id common.BufId) (level int, stack int) {
	_, buf := c.GetActiveBufObj(id)
	if buf == nil {
		return 0, 0
	}
	return buf.GetLevel(), buf.GetStack()
}

func (c *Ctrl) Update() {
	for index := 0; index < len(c.bufs); {
		v := c.bufs[index]
		if !v.IsOver() {
			v.Update()
			index++
		} else {
			c.bufs = slices.Delete(c.bufs, index, index+1)
		}
	}
}

func (c *Ctrl) Destroy() {
	for index := 0; index < len(c.bufs); {
		v := c.bufs[index]
		if !v.IsOver() {
			v.Cancel()
		}
		index++
	}
}
