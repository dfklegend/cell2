package apientry

import (
	"fmt"
	"reflect"
	"strings"

	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/formater"
	"github.com/dfklegend/cell2/utils/logger"
)

const (
	InnerGroupName = "_"
)

// 在启动时集中注册服务
// map就不存在并发问题

func NewCollection() *APICollection {
	collection := &APICollection{
		Containers:  make(map[string]*APIContainer),
		Entries:     &APIEntries{},
		apiFormater: formater.GetDefaultFormater(),
	}
	return collection
}

type APICollection struct {
	// 如果为空，可能没build
	Containers  map[string]*APIContainer
	Entries     *APIEntries
	apiFormater api.IAPIFormatter
}

func (c *APICollection) SetFormater(formater api.IAPIFormatter) {
	c.apiFormater = formater
}

func (c *APICollection) Register(e api.IAPIEntry, options ...Option) *APICollection {
	c.Entries.Register(e, options...)
	return c
}

func (c *APICollection) Build() error {
	// 初始化
	//clear old
	c.Containers = make(map[string]*APIContainer)

	entries := c.Entries.List()
	for _, e := range entries {
		err := c.newService(e.Entry, e.Opts...)
		if err != nil {
			logger.L.Errorf("ExtractHandler failed: %v", err)
			continue
		}
	}

	return nil
}

func (c *APICollection) newService(comp api.IAPIEntry, opts ...Option) error {
	srv := NewContainer(comp, opts...)

	if _, ok := c.Containers[srv.Name]; ok {
		return fmt.Errorf("handler: service already defined: %s", srv.Name)
	}

	if err := srv.ExtractHandler(c.apiFormater); err != nil {
		return err
	}

	// 记录
	c.Containers[srv.Name] = srv
	logger.Log.Infof("newService:%v", srv.Name)
	return nil
}

// Call
/**
 * 1. 转发给对应的Service
 * 2. 按目标参数类型反序列化
 * cb(error, []byte result)
 * 出入口都是[]byte (json序列化)
 * ext
 * 		一个自定义参数，目前是前端接口的session
 *
 *
 * 将请求转发给底层的各接口执行
 * @param route{string} 服务.method路由字符串
 * @param args{[]byte} 请求参数流化的数据(json)
 * @param cbFunc{HandlerCBFunc} 执行完后的回调
 *        func(e error, outArg interface{})
 *        		error 是否执行错误
 *        		outArt json序列化的返回值
 *
 * 注: cb必定返回
 */
func (c *APICollection) Call(ctx api.IContext, route string, arg any,
	cbFunc HandlerCBFunc) {
	serviceName, methodName, goodRoute := c.splitRoute(route)
	if !goodRoute {
		logger.Log.Warnf("bat route: %s", route)
		CheckInvokeCBFunc(cbFunc, fmt.Errorf("bat route: %s", route), nil)
		return
	}

	service := c.Containers[serviceName]
	if service == nil {
		logger.Log.Warnf("can not find service: %s", serviceName)
		CheckInvokeCBFunc(cbFunc, fmt.Errorf("can not find service: %s", serviceName), nil)
		return
	}
	service.CallMethod(ctx, methodName, arg, cbFunc)
}

func (c *APICollection) HasMethod(route string) bool {
	return c.GetArgType(route) != nil
}

func (c *APICollection) splitRoute(route string) (string, string, bool) {
	subs := strings.Split(route, ".")
	if len(subs) == 2 {
		return subs[0], subs[1], true
	}
	if len(subs) == 1 {
		return InnerGroupName, subs[0], true
	}
	//logger.Log.Warnf("bat route: %s", route)
	return "", "", false
}

func (c *APICollection) GetArgType(route string) reflect.Type {
	serviceName, methodName, goodRoute := c.splitRoute(route)
	if !goodRoute {
		return nil
	}

	service := c.Containers[serviceName]
	if service == nil {
		return nil
	}
	return service.GetArgType(methodName)
}

func (c *APICollection) DumpAll() {
	L := logger.L
	for k, v := range c.Containers {
		L.Infof("    apiGroup: %v", k)
		v.Dump()
	}
}
