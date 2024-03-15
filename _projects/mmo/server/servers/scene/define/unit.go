package define

const (
	UnitNone UnitType = iota
	UnitExit          // 出口，场景传送点
	UnitTest
	UnitCamera  // 摄像机
	UnitMonster // 怪物
	UnitAvatar  // 创建出来的玩家化身
)
