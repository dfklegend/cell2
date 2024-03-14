package welcomemodule

import (
	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/baseapp/module"
	"github.com/dfklegend/cell2/utils/logger"
)

type WelcomeModule struct {
	*module.BaseModule
}

func NewWelcomeModule() *WelcomeModule {
	return &WelcomeModule{
		BaseModule: module.NewBaseModule(),
	}
}

func (w *WelcomeModule) Start(next interfaces.FuncWithSucc) {
	logger.Log.Infoln("====================================")
	logger.Log.Infoln("====== Welcome to Cell2 ============")
	logger.Log.Infoln("====================================")
	next(true)
}

func (w *WelcomeModule) Stop(next interfaces.FuncWithSucc) {
	logger.Log.Infoln("====================================")
	logger.Log.Infoln("========= Cell2 Exit... ============")
	logger.Log.Infoln("====================================")
	next(true)
}
