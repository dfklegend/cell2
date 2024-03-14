// 路由模块

package route

import (
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
)

var (
	defaultRouteFunc RouteFunc
)

const (
	BadRouteParam = "bad_route_param"
	MissRouteFunc = "miss_route_func"
	NoService     = "no_service"
)

// IRouteParam
// 抽象，基于session或者map
type IRouteParam interface {
	Get(k string, def interface{}) interface{}
}

type RouteFunc func(serverType string, p IRouteParam) string

func ExamRoute(serverType string, p interface{}) string {
	// . 获取需求的值
	// v := p.Get("examKey").(int)
	// . 从app获取基于服务类型的可选服务器列表
	// . 根据v来获取对应的serverId
	return "some"
}

func SetDefaultRoute(f RouteFunc) {
	defaultRouteFunc = f
}

// --------

type MapParam struct {
	data map[string]interface{}
}

func NewMapParam(data map[string]interface{}) *MapParam {
	return &MapParam{
		data: data,
	}
}

func (p *MapParam) Get(k string, def interface{}) interface{} {
	v, ok := p.data[k]
	if ok {
		return v
	}
	return def
}

var TheRouteService = NewRouteService()

//
type IRouteService interface {
	Register(serverType string, f RouteFunc)
	GetFunc(serverType string) RouteFunc
}

type RouteService struct {
	routes map[string]RouteFunc
}

func NewRouteService() *RouteService {
	return &RouteService{
		routes: make(map[string]RouteFunc),
	}
}

func GetRouteService() *RouteService {
	return TheRouteService
}

func (s *RouteService) Register(serverType string, f RouteFunc) {
	s.routes[serverType] = f
}

func (s *RouteService) GetFunc(serverType string) RouteFunc {
	// 只是读，线程安全
	v, _ := s.routes[serverType]
	return v
}

// Route
// . 如果是nil或者routeparam(frontsession)
// . 如果是map，取值
// . 如果是字符串，直接返回
func (s *RouteService) Route(serverType string, param interface{}) string {
	// 如果是IRouteParam类型
	rparam, ok := param.(IRouteParam)
	if ok || param == nil {
		return s.RouteFromParam(serverType, rparam)
	}

	// 如果是map
	m, ok := param.(map[string]interface{})
	if ok {
		return s.RouteFromMap(serverType, m)
	}

	// 如果是字符串
	str, ok := param.(string)
	if ok {
		return str
	}

	l.Log.Errorf("error route param:%v\n", param)
	return BadRouteParam
}

func (s *RouteService) doRoute(serverType string, param IRouteParam) string {
	f := s.GetFunc(serverType)
	if f == nil && defaultRouteFunc != nil {
		f = defaultRouteFunc
	}

	if f == nil {
		return MissRouteFunc
	}

	defer func() {
		if err := recover(); err != nil {
			l.E.Errorf("panic in Route:%v", err)
			l.E.Errorf(common.GetStackStr())
		}
	}()

	return f(serverType, param)
}

func (s *RouteService) RouteFromParam(serverType string, param IRouteParam) string {
	return s.doRoute(serverType, param)
}

func (s *RouteService) RouteFromMap(serverType string, m map[string]interface{}) string {
	return s.doRoute(serverType, NewMapParam(m))
}
