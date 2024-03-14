package impls

import (
	nodeapi2 "github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/config"
	"github.com/dfklegend/cell2/node/service"
)

//	ServiceCreateHandler 后端服务器
func ServiceCreateHandler(service *service.NodeService, info *config.ServiceInfo) {
	handler := NewHandler(nodeapi2.Registry.GetCollection(app.MakeName(info.Type, "handler")))

	service.AddComponent("handler", handler)
}

func ServiceCreateCommonComponents(service *service.NodeService, info *config.ServiceInfo) {
	service.AddComponent("channel", NewChannelComponent())
	ServiceCreateHandler(service, info)
}
