package base

import (
	"github.com/dfklegend/cell2/utils/csvutils"
	l "github.com/dfklegend/cell2/utils/logger"
)

type IStringKeyEntry interface {
	GetId() string
}

// DataMgr
// 通过字符串的key来查找
type DataMgr[T IStringKeyEntry] struct {
	buf     []T
	entries map[string]T
}

func NewDataMgr[T IStringKeyEntry]() *DataMgr[T] {
	return &DataMgr[T]{
		buf:     []T{},
		entries: map[string]T{},
	}
}

func (m *DataMgr[T]) LoadFromFile(path string) error {
	err := csvutils.LoadFromFile(path, &m.buf)
	if err != nil {
		l.L.Errorf("csv LoadFromFile: %v, error: %v", path, err)
		return err
	}

	for _, v := range m.buf {
		key := v.GetId()
		if key == "" {
			continue
		}
		m.entries[key] = v
	}
	return nil
}

func (m *DataMgr[T]) GetEntry(id string) T {
	return m.entries[id]
}

// AddEntry for test
func (m *DataMgr[T]) AddEntry(entry T) {
	m.buf = append(m.buf, entry)
	key := entry.GetId()
	if key != "" {
		m.entries[key] = entry
	}
}
