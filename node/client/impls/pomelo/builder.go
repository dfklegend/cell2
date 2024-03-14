package pomelo

import (
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/client/impls"
	"github.com/dfklegend/cell2/node/config"
	"github.com/dfklegend/cell2/node/service"
)

func ServiceCreateAcceptors(service *service.NodeService, name string, info *config.ServiceInfo) {
	if info == nil || !info.Frontend {
		return
	}
	if info.ClientAddress == "" && info.WSClientAddress == "" {
		return
	}

	sessions := impls.NewClientSessions(name)

	var handler *impls.HandlerComponent
	handler, _ = service.GetComponent("handler").(*impls.HandlerComponent)
	if handler == nil {
		handler = impls.NewHandler(registry.Registry.GetCollection(app.MakeName(info.Type, "handler")))
		service.AddComponent("handler", handler)
	}
	forwarder := impls.NewForwarder()

	service.AddComponent("sessions", impls.NewSessionsComponent(sessions))
	service.AddComponent("forwarder", forwarder)

	sessions.SetHandler(handler)
	handler.SetForwarder(forwarder)

	// acceptors
	// pomelo-tcp
	if info.ClientAddress != "" {
		tcp := NewTCPComponent(sessions)
		service.AddComponent("tcp", tcp)
		tcp.Start(info.ClientAddress)
	}

	// pomelo-ws
	if info.WSClientAddress != "" {
		ws := NewWSComponent(sessions)
		service.AddComponent("ws", ws)
		ws.Start(info.WSClientAddress)
	}
}
