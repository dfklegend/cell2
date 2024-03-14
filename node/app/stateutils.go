package app

import (
	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
)

func NotifyServiceRetired(ns *service.NodeService) {
	ns.RequestEx(Node.GetNodeCtrl().GetAdmin(), "ctrl.servicecmd", &msgs.ServiceCmd{
		Name: ns.Name,
		Cmd:  "retired",
	}, func(err error, raw interface{}) {
		if err != nil {
			return
		}
		ack, _ := raw.(*msgs.ServiceCmdAck)
		if ack == nil {
			return
		}
		l.L.Infof("notify retired, ack: %v", ack.Result)
	})
}
