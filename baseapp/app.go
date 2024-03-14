package baseapp

import (
	"log"

	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/baseapp/module"
	"github.com/dfklegend/cell2/utils/runservice"
)

type App struct {
	rs    *runservice.StandardRunService
	ms    *module.ModList
	state int
}

func NewApp() *App {
	return &App{
		ms:    module.NewModList(),
		state: State0,
	}
}

func (a *App) setState(s int) {
	a.state = s
}

func (a *App) isState(s int) bool {
	return a.state == s
}

func (a *App) GetRunService() *runservice.StandardRunService {
	return a.rs
}

func (a *App) Prepare() {
	a.rs = runservice.NewStandardRunService("__App__")
	a.rs.Start()

	a.setState(StatePrepared)
}

func (a *App) Cleanup() {
	a.rs.Stop()
	log.Printf("App.Cleanup")
}

func (a *App) AddModule(mod interfaces.IAppModule) {
	a.ms.AddModule(mod)
	mod.Init(a.rs)
}

func (a *App) Start(finish interfaces.FuncWithSucc) {
	if !a.isState(StatePrepared) {
		return
	}

	a.setState(StateStarting)
	a.ms.Start(func(succ bool) {
		if succ {
			a.setState(StateNormal)
		}
		if finish != nil {
			finish(succ)
		}
	})
}

func (a *App) Stop(finish interfaces.FuncWithSucc) {
	if !a.isState(StateNormal) {
		return
	}

	a.setState(StateStoping)
	a.ms.Stop(func(succ bool) {
		if succ {
			a.setState(StateStopped)
		}
		if finish != nil {
			finish(succ)
		}

		if succ {
			a.Cleanup()
		}
	})
}
