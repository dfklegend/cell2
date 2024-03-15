package impls

// AttrSlice 属性集合
type AttrSlice struct {
	size  int
	attrs []Attr
}

func NewAttrSlice(size int) *AttrSlice {
	return &AttrSlice{
		size:  size,
		attrs: make([]Attr, size, size),
	}
}

func (a *AttrSlice) Clone() *AttrSlice {
	v1 := NewAttrSlice(a.size)
	for k, v := range a.attrs {
		v.Copy(&v1.attrs[k])
	}
	return v1
}

func (a *AttrSlice) GetAttr(index int) *Attr {
	if index < 0 || index >= a.size {
		return nil
	}
	return &a.attrs[index]
}

func (a *AttrSlice) GetNum() int {
	return len(a.attrs)
}

func (a *AttrSlice) Reset() {
	for i := 0; i < len(a.attrs); i++ {
		a.attrs[i].Reset()
	}
}
