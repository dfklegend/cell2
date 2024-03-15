package main

import (
	"github.com/dfklegend/cell2/node/route"
)

func initRoutes() {
	rs := route.GetRouteService()
	registerRoutes(rs)
}

func registerRoutes(rs *route.RouteService) {
	// logic_old
	rs.Register("logic", func(serviceType string, param route.IRouteParam) string {
		id := param.Get("logic", "").(string)
		return id
	})

	rs.Register("scene", func(serviceType string, param route.IRouteParam) string {
		id := param.Get("scene", "").(string)
		return id
	})
}
