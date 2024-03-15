package ai

// 能移动的状态机
// AIState
const (
	StateInit         = iota
	StateWait         // 等待
	StateRandMove     // 随机移动
	StateAttack       // 攻击
	StateNoMoveAttack // 不会移动，只攻击在范围内的敌人
	StateDead         // 死亡状态
)

// 卡牌的state
const (
	CardStateWait = iota + 100
	CardStateAttack
	CardStateDead
)
