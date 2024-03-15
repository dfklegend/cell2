package attrsync

import (
	"mmo/messages/cproto"
)

type ISyncTarget interface {
	Push(attrs *cproto.AttrsChanged)
}

// Syncer 属性同步器
type Syncer struct {
	attrs *cproto.AttrsChanged
}

func NewSyncer() *Syncer {
	return &Syncer{
		attrs: &cproto.AttrsChanged{},
	}
}

func (s *Syncer) Clear() {

}

func (s *Syncer) SetId(id int32) {
	s.attrs.Id = id
}

func (s *Syncer) AddIntValue(index int32, value int32) {
	attrs := s.attrs.Ints
	if attrs == nil {
		attrs = make([]*cproto.IntAttr, 0)
		s.attrs.Ints = attrs
	}

	attrs = append(attrs, &cproto.IntAttr{
		Index: index,
		Value: value,
	})
}

func (s *Syncer) AddFloatValue(index int32, value float32) {
	attrs := s.attrs.Floats
	if attrs == nil {
		attrs = make([]*cproto.FloatAttr, 0)
		s.attrs.Floats = attrs
	}
	attrs = append(attrs, &cproto.FloatAttr{
		Index: index,
		Value: value,
	})
}

func (s *Syncer) Push(tar ISyncTarget) {
	tar.Push(s.attrs)
}
