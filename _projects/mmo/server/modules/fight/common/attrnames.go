package common

var attrNames = newAttrNames()

func RegisterAttr(index int, name string) {
	attrNames.Add(index, name)
}

func AttrNameToIndex(name string) int {
	return attrNames.NameToIndex(name)
}

func AttrIndexToName(index int) string {
	return attrNames.IndexToName(index)
}

type AttrNames struct {
	indexToName map[int]string
	nameToIndex map[string]int
}

func newAttrNames() *AttrNames {
	return &AttrNames{
		indexToName: map[int]string{},
		nameToIndex: map[string]int{},
	}
}

func (a *AttrNames) Add(index int, name string) {
	a.indexToName[index] = name
	a.nameToIndex[name] = index
}

func (a *AttrNames) NameToIndex(name string) int {
	index, ok := a.nameToIndex[name]
	if !ok {
		return -1
	}
	return index
}

func (a *AttrNames) IndexToName(index int) string {
	name, ok := a.indexToName[index]
	if !ok {
		return "invalid attr"
	}
	return name
}
