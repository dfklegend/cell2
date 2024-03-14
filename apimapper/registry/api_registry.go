package registry

import (
	"sync"

	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/utils/logger"
)

/*
	全局可以注册APICollection，并查询使用
*/

var (
	Registry = newRegistry()
)

// 	APIRegistry
//	接口注册
type APIRegistry struct {
	apis map[string]*apientry.APICollection
	sync.RWMutex
}

func newRegistry() *APIRegistry {
	return &APIRegistry{
		apis: make(map[string]*apientry.APICollection),
	}
}

func (r *APIRegistry) AddCollection(name string) *apientry.APICollection {
	r.Lock()
	defer r.Unlock()

	item := r.apis[name]
	if item != nil {
		return item
	}
	item = apientry.NewCollection()
	r.apis[name] = item
	return item
}

func (r *APIRegistry) GetCollection(name string) *apientry.APICollection {
	r.RLock()
	defer r.RUnlock()

	return r.apis[name]
}

func (r *APIRegistry) Build() {
	r.RLock()
	defer r.RUnlock()

	for _, v := range r.apis {
		v.Build()
	}
}

//	为了方便调试
// 	输出所有的API
func (r *APIRegistry) DumpAll() {
	r.RLock()
	defer r.RUnlock()

	L := logger.L
	L.Infof(" ---- dump all apis ---- ")
	for k, v := range r.apis {
		L.Infof("  collection: %v", k)
		v.DumpAll()
	}
	L.Infof(" ---- dump over ---- ")
}
