package define

/*
	处理节点控制命令
		stat
		retire
		exit
*/

// NodeState 节点状态
type NodeState int

const (
	Init NodeState = iota
	Working
	Retiring
	Retired
	Exiting
	Exited
)

const (
	NodeAdmin = "__nodeadmin__"

	// TODO: 命令都写成常量
)
