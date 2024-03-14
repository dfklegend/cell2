package app

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/node/route"
	"github.com/dfklegend/cell2/nodectrl/define"
	l "github.com/dfklegend/cell2/utils/logger"
)

/*
	可以判断服务处于什么状态
	常规的route不会管service处于什么状态
	一般，负载均衡模块，分配时，需要主动获取处于work状态的service
	比如: GetFirstWorkService, RandGetWorkService
*/

func init() {
	route.SetDefaultRoute(defaultRoute)
}

func IsWorkState(state int) bool {
	return state == int(define.Working)
}

func GetServicePID(name string) *actor.PID {
	s := Node.GetCluster().GetService(name)
	if s == nil {
		return nil
	}
	return s.PID
}

func GetWorkServicePID(name string) *actor.PID {
	s := Node.GetCluster().GetService(name)
	if s == nil || !IsWorkState(s.State) {
		return nil
	}
	return s.PID
}

func GetFirstService(serviceType string) *actor.PID {
	nodes := Node.GetCluster().GetServiceList(serviceType)
	if nodes == nil || len(nodes.Items) == 0 {
		return nil
	}
	return nodes.Items[0].PID
}

func GetFirstWorkService(serviceType string) *actor.PID {
	nodes := Node.GetCluster().GetWorkServiceList(serviceType)
	if nodes == nil || len(nodes.Items) == 0 {
		return nil
	}
	return nodes.Items[0].PID
}

func GetFirstServiceItem(serviceType string) *ServiceItem {
	nodes := Node.GetCluster().GetServiceList(serviceType)
	if nodes == nil || len(nodes.Items) == 0 {
		return nil
	}
	return nodes.Items[0]
}

func GetFirstWorkServiceItem(serviceType string) *ServiceItem {
	nodes := Node.GetCluster().GetWorkServiceList(serviceType)
	if nodes == nil || len(nodes.Items) == 0 {
		return nil
	}
	return nodes.Items[0]
}

func RandGetService(serviceType string) *actor.PID {
	item := RandGetServiceItem(serviceType)
	if item == nil {
		return nil
	}
	return item.PID
}

func RandGetWorkService(serviceType string) *actor.PID {
	item := RandGetWorkServiceItem(serviceType)
	if item == nil {
		return nil
	}
	return item.PID
}

func RandGetServiceItem(serviceType string) *ServiceItem {
	nodes := Node.GetCluster().GetServiceList(serviceType)
	if nodes == nil || len(nodes.Items) == 0 {
		return nil
	}
	index := rand.Intn(len(nodes.Items))
	return nodes.Items[index]
}

func RandGetWorkServiceItem(serviceType string) *ServiceItem {
	nodes := Node.GetCluster().GetWorkServiceList(serviceType)
	if nodes == nil || len(nodes.Items) == 0 {
		return nil
	}
	index := rand.Intn(len(nodes.Items))
	return nodes.Items[index]
}

func RandGetServiceName(serviceType string) string {
	item := RandGetServiceItem(serviceType)
	if item == nil {
		return ""
	}
	return item.Name
}

func RandGetWorkServiceName(serviceType string) string {
	item := RandGetWorkServiceItem(serviceType)
	if item == nil {
		return ""
	}
	return item.Name
}

func GetServices(serviceType string) *ServiceList {
	return Node.GetCluster().GetServiceList(serviceType)
}

func GetWorkServices(serviceType string) *ServiceList {
	return Node.GetCluster().GetWorkServiceList(serviceType)
}

// 	serviceType.category.method
func SplitClientRoute(route string) (string, string, string) {
	subs := strings.Split(route, ".")
	if len(subs) != 3 {
		return "", "", ""
	}
	return subs[0], subs[1], subs[2]
}

func RoutePID(serviceType string, param any) *actor.PID {
	id := route.GetRouteService().Route(serviceType, param)
	if id == "" {
		return nil
	}
	return GetServicePID(id)
}

// select first
func defaultRoute(serverType string, param route.IRouteParam) string {
	servers := GetWorkServices(serverType)
	if servers == nil || len(servers.Items) == 0 {
		l.Log.Errorf("can not find serverType:%v", serverType)
		return route.NoService
	}

	return servers.Items[0].Name
}

//	xx.remote rpc接口
// 	xx.handler 前端接口
func MakeName(serviceType string, postfix string) string {
	return fmt.Sprintf("%v.%v", serviceType, postfix)
}

func splitHostPort(addr string) (host string, port int, err error) {
	if h, p, e := net.SplitHostPort(addr); e != nil {
		if addr != "nonhost" {
			err = e
		}
		host = "nonhost"
		port = -1
	} else {
		host = h
		port, err = strconv.Atoi(p)
	}
	return
}
