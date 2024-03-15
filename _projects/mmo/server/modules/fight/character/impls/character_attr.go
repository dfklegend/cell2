package charimpls

import (
	"golang.org/x/exp/slices"

	"mmo/modules/fight/attr"
	"mmo/modules/fight/common"
)

func (c *Character) InitAttrWatchers() {
	c.GetAttr(common.AttackSpeed).SetWatcher(attr.NewFuncWatcher(func(attr attr.IAttr) {
		c.onAttackSpeedDirt(attr)
	}))
}

func (c *Character) VisitAttr(doFunc func(index int, item attr.IAttr)) {
	num := c.attrs.GetNum()
	for i := 0; i < num; i++ {
		item := c.attrs.GetAttr(i)
		doFunc(i, item)
	}
}

// 根据攻速，计算普攻间隔
func (c *Character) onAttackSpeedDirt(attr attr.IAttr) {
	speed := attr.GetValue()
	if speed <= 0 {
		return
	}
	c.skillTable.SetNormalAttackInterval(float32(1.0 / speed))
}

func (c *Character) RecalculateAttrs() {
	c.attrs.Reset()
	// 重新挂一下
	for i := 0; i < len(c.equipGroups); i++ {
		c.equipGroups[i].Equip(c)
	}
}

func (c *Character) OffsetBase(index int, value attr.Value) {
	//c.GetAttr(index).OffsetBase(value)
	c.applyAttr(index, func(a attr.IAttr) {
		a.OffsetBase(value)
	})
}

func (c *Character) OffsetIntBase(index int, value int) {
	//c.GetAttr(index).OffsetBase(attr.Value(value))
	c.applyAttr(index, func(a attr.IAttr) {
		a.OffsetBase(attr.Value(value))
	})
}

func (c *Character) OffsetPercent(index int, off float32) {
	//c.GetAttr(index).OffsetPercent(off)
	c.applyAttr(index, func(a attr.IAttr) {
		a.OffsetPercent(off)
	})
}

func (c *Character) SetBaseValue(index int, v attr.Value) {
	//c.GetAttr(index).SetBase(v)
	c.applyAttr(index, func(a attr.IAttr) {
		a.SetBase(v)
	})
}

func (c *Character) SetIntBaseValue(index int, v int) {
	//c.GetAttr(index).SetBase(attr.Value(v))
	c.applyAttr(index, func(a attr.IAttr) {
		a.SetBase(attr.Value(v))
	})
}

func (c *Character) applyAttr(index int, apply func(a attr.IAttr)) {
	a := c.GetAttr(index)
	oldV := a.GetValue()
	apply(a)
	newV := a.GetValue()

	off := newV - oldV
	if off > 0.00001 || off < -0.00001 {
		c.onAttrChanged(index, oldV, newV)
	}
}

func (c *Character) GetValue(index int) attr.Value {
	return c.GetAttr(index).GetValue()
}

func (c *Character) SetIntValue(index int, v int) {
	c.SetIntBaseValue(index, v)
}

func (c *Character) GetIntValue(index int) int {
	return c.GetAttr(index).GetIntValue()
}

func (c *Character) onAttrChanged(index int, old, new attr.Value) {
	//l.L.Infof(" %v attrChanged %v %v -> %v", c.id, common.AttrIndexToName(index), old, new)
	for _, v := range c.attrCWatchers {
		v.OnAttrChanged(index, old, new)
	}
}

func (c *Character) setLevel(level int) {
	c.GetAttr(common.Level).SetBase(attr.Value(level))
}

func (c *Character) GetLevel() int {
	return c.GetIntValue(common.Level)
}

func (c *Character) ChangeLevel(level int) {
	old := c.GetLevel()
	if old == level {
		return
	}
	c.setLevel(level)
	c.events.Publish("onlevelchanged", old, level)
}

func (c *Character) GetSide() int {
	return c.GetIntValue(common.Side)
}

func (c *Character) SetHP(hp int) {
	hpMax := c.GetIntValue(common.HPMax)
	if hp > hpMax {
		hp = hpMax
	}
	c.SetIntBaseValue(common.HP, hp)
}

func (c *Character) GetHP() int {
	return c.GetIntValue(common.HP)
}

func (c *Character) SetEnergy(v int) {
	c.SetIntValue(common.Energy, v)
}

func (c *Character) GetEnergy() int {
	return c.GetIntValue(common.Energy)
}

func (c *Character) AddEnergy(off int) int {
	c.OffsetIntBase(common.Energy, off)
	v := c.GetEnergy()
	max := c.GetIntValue(common.EnergyMax)
	if v > max {
		c.SetEnergy(max)
		return max
	}
	if v < 0 {
		v = 0
	}
	return v
}

func (c *Character) ApplyMPCost(mpcost int) {
	c.addEnergyAndBroadcast(-mpcost)
}

func (c *Character) GetNormalAttackInterval() int32 {
	return c.normalAttackInterval
}

func (c *Character) AddAttrChangeWatcher(watcher common.IAttrChangeWatcher) {
	c.attrCWatchers = append(c.attrCWatchers, watcher)
}

// RemoveAttrChangeWatcher TODO: 多watcher
func (c *Character) RemoveAttrChangeWatcher(watcher common.IAttrChangeWatcher) {
	//index := slices.Index(c.attrCWatchers, watcher)
	//if index == -1 {
	//	return
	//}
	//c.attrCWatchers = slices.Delete(c.attrCWatchers, index, index+1)
	if len(c.attrCWatchers) == 0 {
		return
	}
	c.attrCWatchers = slices.Delete(c.attrCWatchers, 0, 1)
}

func (c *Character) broadcastAttrValue(which int) {
	if c.world.GetWatcher() == nil {
		return
	}
	data := &common.DataAttrsChanged{
		Attrs: make([]*common.OneAttr, 0),
	}

	data.Attrs = append(data.Attrs, &common.OneAttr{
		Index: which,
		Value: c.GetValue(which),
	})

	c.world.GetWatcher().OnAttrsChanged(c, data)
}
