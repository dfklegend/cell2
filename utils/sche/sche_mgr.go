package sche

import (
	"sync"
)

// 	ScheMgr --------------------------
// 	管理器,通过名字来创建和访问Sche
type Mgr struct {
	sches map[string]*Sche
	mutex sync.Mutex
}

func NewScheMgr() *Mgr {
	return &Mgr{
		sches: make(map[string]*Sche),
		//mutex:
	}
}

// GetSche createIfMiss
func (s *Mgr) GetSche(name string) *Sche {
	defer s.mutex.Unlock()

	s.mutex.Lock()
	one := s.sches[name]
	if one == nil {
		s.sches[name] = NewSche()
		return s.sches[name]
	}
	return one
}

func (s *Mgr) hasSche(name string) bool {
	defer s.mutex.Unlock()

	s.mutex.Lock()
	one := s.sches[name]
	return one != nil
}

func (s *Mgr) DelSche(name string) {
	defer s.mutex.Unlock()

	s.mutex.Lock()
	delete(s.sches, name)
}
