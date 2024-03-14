package builder

import (
	"testing"

	"github.com/dfklegend/cell2/baseapp"
	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/node/route"
	"github.com/dfklegend/cell2/node/service"
	utils "github.com/dfklegend/cell2/nodeutils"
	l "github.com/dfklegend/cell2/utils/logger"
)

func TestBuilder(t *testing.T) {
	nodeId := "node-1"
	NewBuilder().
		ConfigLog(nodeId, "./logs").
		SetDefaultLaunchFunc(func(app interfaces.IApp) {
			l.L.Infof("default mode")
			utils.NodeAddCommonModules(app)
			utils.NodeAddClusterModules(app)
		}).
		RegisterLaunchModes(func() {
			baseapp.LaunchFunc("allinone", func(app interfaces.IApp) {
				l.L.Infof("allinone mode")
				utils.NodeAddCommonModules(app)
				utils.NodeAddClusterModules(app)
			})

			baseapp.LaunchFunc("gate", func(app interfaces.IApp) {
				l.L.Infof("gate mode")
				utils.NodeAddCommonModules(app)
				utils.NodeAddClusterModules(app)
			})
		}).
		RegisterAPIs(true, func() {
			//
		}).
		RegisterServiceCreators(func(factory service.IServiceFactory) {
			//
		}).
		RegisterRoutes(func(rs *route.RouteService) {
			//
		}).
		Next(func() {

		}).
		StartApp("../testdata/config", nodeId, func(succ bool) {
			l.L.Infof("startApp %v", succ)
		})
}
