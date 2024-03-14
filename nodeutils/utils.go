package utils

/*
	避免交叉引用
*/

import (
	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/node/builtin"

	actormodule "github.com/dfklegend/cell2/node/modules/actor"
	clustermodule "github.com/dfklegend/cell2/node/modules/cluster"
	welcomemodule "github.com/dfklegend/cell2/node/modules/welcome"
)

func NodeAddCommonModules(app interfaces.IApp) {
	app.AddModule(welcomemodule.NewWelcomeModule())
	app.AddModule(actormodule.NewActorSystemModule())
}

func NodeAddClusterModules(app interfaces.IApp) {
	app.AddModule(clustermodule.NewClusterModule())
}

func NodeInitSystemAPI() {
	builtin.Visit()
}
