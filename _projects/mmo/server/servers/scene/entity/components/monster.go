package components

type MonsterComponent struct {
	*BaseSceneComponent
	cfgId string
}

func NewMonsterComponent(cfgId string) *MonsterComponent {
	return &MonsterComponent{
		BaseSceneComponent: NewBaseSceneComponent(),
		cfgId:              cfgId,
	}
}

func (c *MonsterComponent) OnPrepare() {
	c.BaseSceneComponent.OnPrepare()
}

func (c *MonsterComponent) OnStart() {
}

func (c *MonsterComponent) Update() {
}

func (c *MonsterComponent) LateUpdate() {
}

func (c *MonsterComponent) OnDestroy() {
}

func (c *MonsterComponent) GetCfgId() string {
	return c.cfgId
}
