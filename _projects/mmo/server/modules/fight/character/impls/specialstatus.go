package charimpls

import (
	"mmo/modules/fight/common"
)

type SpecialStatus struct {
	Id  int
	Ref int
}

type SpecialStatusCtrl struct {
	owner  common.ICharacter
	status []*SpecialStatus
}

func NewSpecialStatusCtrl(owner common.ICharacter) *SpecialStatusCtrl {
	c := &SpecialStatusCtrl{
		owner: owner,
	}
	c.init(common.SSMax)
	return c
}

func (c *SpecialStatusCtrl) init(size int) {

	c.status = make([]*SpecialStatus, size)
	for k, _ := range c.status {
		c.status[k] = &SpecialStatus{
			Id: k,
		}
	}
}

func (c *SpecialStatusCtrl) getItem(id int) *SpecialStatus {
	if id >= common.SSMax {
		return nil
	}
	return c.status[id]
}

func (c *SpecialStatusCtrl) offSpecialStatus(id int, off int) {
	s := c.getItem(id)
	if s == nil {
		return
	}
	oldV := s.Ref > 0
	s.Ref += off
	newV := s.Ref > 0

	if oldV != newV {
		c.owner.OnSpecialStatusChanged(id, oldV, newV)
	}

}

func (c *SpecialStatusCtrl) AddSpecialStatus(id int) {
	c.offSpecialStatus(id, 1)
}

func (c *SpecialStatusCtrl) SubSpecialStatus(id int) {
	c.offSpecialStatus(id, -1)
}

func (c *SpecialStatusCtrl) HasSpecialStatus(id int) bool {
	s := c.getItem(id)
	if s == nil {
		return false
	}
	return s.Ref > 0
}
