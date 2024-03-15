package entry

// MonsterTemplate
// 怪物模板数据
type MonsterTemplate struct {
	Id      string  `csv:"id"`
	HP      float32 `csv:"hp"`
	HPLv    float32 `csv:"hpLv"`
	Armor   float32 `csv:"armor"`
	ArmorLv float32 `csv:"armorLv"`
}

func (s *MonsterTemplate) GetId() string {
	return s.Id
}

// Monster
// 怪物数据
type Monster struct {
	Id       string `csv:"id"`
	Name     string `csv:"name"`
	Template string `csv:"template"`
	Level    int    `csv:"level"`
}

func (s *Monster) GetId() string {
	return s.Id
}
