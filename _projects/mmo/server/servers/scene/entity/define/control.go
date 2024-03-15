package define

type IControl interface {
	ReqMoveTo(x, z float32)
	ReqStopMove(x, z float32)
}
