package service

import (
	"errors"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"

	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/utils/serialize"
	"github.com/dfklegend/cell2/utils/serialize/proto"
)

// APIDispatcher
// 将消息转发给目标接口
type APIDispatcher struct {
	apis []*apientry.APICollection
	ctx  *RemoteContext
}

func NewDispatcher(cols ...*apientry.APICollection) *APIDispatcher {
	d := &APIDispatcher{
		apis: make([]*apientry.APICollection, 0),
		ctx:  NewRemoteContext(),
	}

	if cols != nil {
		d.AddCollection(cols...)
	}
	return d
}

func (d *APIDispatcher) AddCollection(cols ...*apientry.APICollection) {
	if cols != nil {
		for _, v := range cols {
			if v == nil {
				continue
			}
			d.apis = append(d.apis, v)
		}
	}
}

func (d *APIDispatcher) Dispatch(ctx actor.Context, service *Service,
	route string, request *messages.ServiceRequest) bool {

	d.ctx.Update(ctx)

	serializer := proto.GetDefaultSerializer()

	e := d.tryCall(service, route, request, serializer)
	// 直接返回错误
	if e != nil {
		service.Response(request, CodeErrString, e.Error(), nil)
	}
	return e == nil
}

func (d *APIDispatcher) tryCall(service *Service, route string,
	request *messages.ServiceRequest, serializer serialize.Serializer) error {
	for i := 0; i < len(d.apis); i++ {
		col := d.apis[i]

		processed := d.tryCallCol(service, col, route, request, serializer)
		if !processed {
			continue
		}
		// succ
		return nil
	}
	return errors.New(fmt.Sprintf("no method: %v", route))
}

// 	return
//		if process by this collection
func (d *APIDispatcher) tryCallCol(service *Service, col *apientry.APICollection, route string,
	request *messages.ServiceRequest, serializer serialize.Serializer) bool {
	if !col.HasMethod(route) {
		return false
	}

	if request.ReqId == NotifyReqID {
		// notify
		apientry.CallWithSerialize(col, d.ctx, route, request.Body, nil, serializer)
		return true
	}
	apientry.CallWithSerialize(col, d.ctx, route, request.Body, func(err error, ret interface{}) {
		if err != nil {
			service.Response(request, CodeErrString, err.Error(), nil)
			return
		}
		service.Response(request, CodeSucc, "", ret)
	}, serializer)
	return true
}
