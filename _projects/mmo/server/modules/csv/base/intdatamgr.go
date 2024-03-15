package base

import (
	"github.com/dfklegend/cell2/utils/csvutils"
	l "github.com/dfklegend/cell2/utils/logger"
)

type IIntKeyEntry interface {
	GetId() int
}

// IntDataMgr
// 通过字符串的key来查找
type IntDataMgr[T IIntKeyEntry] struct {
	buf     []T
	entries map[int]T
}

func NewIntDataMgr[T IIntKeyEntry]() *IntDataMgr[T] {
	return &IntDataMgr[T]{
		buf:     []T{},
		entries: map[int]T{},
	}
}

func (m *IntDataMgr[T]) LoadFromFile(path string) error {
	err := csvutils.LoadFromFile(path, &m.buf)
	if err != nil {
		l.L.Errorf("csv LoadFromFile: %v, error: %v", path, err)
		return err
	}

	for _, v := range m.buf {
		key := v.GetId()
		if key == 0 {
			continue
		}
		m.entries[key] = v
	}
	return nil
}

func (m *IntDataMgr[T]) GetEntry(id int) T {
	return m.entries[id]
}

func (m *IntDataMgr[T]) Visit(visitor func(T)) {
	for _, v := range m.entries {
		visitor(v)
	}
}

// AddEntry for test
func (m *IntDataMgr[T]) AddEntry(entry T) {
	m.buf = append(m.buf, entry)
	key := entry.GetId()
	if key != 0 {
		m.entries[key] = entry
	}
}
