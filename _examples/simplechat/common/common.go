package common

import (
	"math/rand"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/node/app"
)

func GetFirstService(serviceType string) *actor.PID {
	nodes := app.Node.GetCluster().GetServiceList(serviceType)
	if nodes == nil || len(nodes.Items) == 0 {
		return nil
	}
	return nodes.Items[0].PID
}

func RandGetService(serviceType string) *actor.PID {
	nodes := app.Node.GetCluster().GetServiceList(serviceType)
	if nodes == nil || len(nodes.Items) == 0 {
		return nil
	}
	index := rand.Intn(len(nodes.Items))
	return nodes.Items[index].PID
}

func RandGetServiceItem(serviceType string) *app.ServiceItem {
	nodes := app.Node.GetCluster().GetServiceList(serviceType)
	if nodes == nil || len(nodes.Items) == 0 {
		return nil
	}
	index := rand.Intn(len(nodes.Items))
	return nodes.Items[index]
}
