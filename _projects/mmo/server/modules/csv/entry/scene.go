package entry

type Scene struct {
	Id            int    `csv:"id"`
	Name          string `csv:"name"`
	Type          int32  `csv:"type"`
	Width         int32  `csv:"width"`
	Height        int32  `csv:"height"`
	MaxMonsterNum int32  `csv:"maxMonsterNum"`
}

func (s *Scene) GetId() int {
	return s.Id
}
