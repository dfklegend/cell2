package cmd

var (
	// TheMgr 方便使用
	TheMgr *Mgr
)

func init() {
	tryCreate()
}

func Visit() {
	tryCreate()
}

func tryCreate() {
	if TheMgr != nil {
		return
	}
	TheMgr = NewCmdMgr()
}

func RegisterCmd(c ICmd) {
	TheMgr.Register(c)
}

func RegisterFuncCmd(name string, cb func(args []string)) {
	RegisterCmd(NewFuncCmd(name, cb))
}

func DispatchCmd(content string) {
	TheMgr.Dispatch(content)
}
