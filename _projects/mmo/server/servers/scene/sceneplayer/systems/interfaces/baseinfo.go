package interfaces

type IBaseInfo interface {
	GetLoginTimes() int32

	AddExp(exp int)
	GetLevel() int32
	AddMoney(v int)

	PushSystemInfo(t int32, info string)
	PushInfoUpdate()
}
