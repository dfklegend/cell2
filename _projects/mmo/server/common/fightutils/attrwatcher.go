package fightutils

import (
	"golang.org/x/exp/slices"

	"mmo/modules/fight/attr"
)

type Attr struct {
	Index int
	OldV  attr.Value
	NewV  attr.Value
}

// AttrChangeWatcher
// 监视属性变化
type AttrChangeWatcher struct {
	attrs []*Attr
}

func NewAttrWatcher() *AttrChangeWatcher {
	return &AttrChangeWatcher{
		attrs: make([]*Attr, 0, 10),
	}
}

func (w *AttrChangeWatcher) Reset() {
	w.attrs = make([]*Attr, 0, 10)
}

func (w *AttrChangeWatcher) findIndex(attrIndex int) int {
	return slices.IndexFunc(w.attrs, func(one *Attr) bool {
		return attrIndex == one.Index
	})
}

func (w *AttrChangeWatcher) OnAttrChanged(attrIndex int, oldV, newV attr.Value) {
	index := w.findIndex(attrIndex)
	if index == -1 {
		one := &Attr{
			Index: attrIndex,
			OldV:  oldV,
			NewV:  newV,
		}
		w.attrs = append(w.attrs, one)
		return
	}

	one := w.attrs[index]
	one.NewV = newV
}

func (w *AttrChangeWatcher) Visit(doFunc func(one *Attr)) {
	for _, v := range w.attrs {
		doFunc(v)
	}
}
