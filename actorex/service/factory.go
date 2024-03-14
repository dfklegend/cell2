package service

import (
	"strings"

	"github.com/asynkron/protoactor-go/actor"

	dispp "github.com/dfklegend/cell2/actorex/disp"
	"github.com/dfklegend/cell2/actorex/mailbox"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	l "github.com/dfklegend/cell2/utils/logger"
)

//	一些扩展的props
type ExtProps struct {
	PostFuncs  []func(service IService)
	APIs       string
	Dispatcher IAPIDispatcher

	PostStartFuncs []func(ctx actor.Context)
}

func (e *ExtProps) WithPostFunc(f func(service IService)) *ExtProps {
	if e.PostFuncs == nil {
		e.PostFuncs = make([]func(service IService), 0)
	}
	e.PostFuncs = append(e.PostFuncs, f)
	return e
}

func (e *ExtProps) WithPostStartFunc(f func(ctx actor.Context)) *ExtProps {
	if e.PostStartFuncs == nil {
		e.PostStartFuncs = make([]func(ctx actor.Context), 0)
	}
	e.PostStartFuncs = append(e.PostStartFuncs, f)
	return e
}

// WithDispatcher 消息分派器
func (e *ExtProps) WithDispatcher(d IAPIDispatcher) *ExtProps {
	e.Dispatcher = d
	return e
}

func (e *ExtProps) WithAPIs(apis string) *ExtProps {
	e.APIs = apis
	return e
}

func (e *ExtProps) doPostFuncs(service IService) {
	if e.PostFuncs == nil {
		return
	}
	for i := 0; i < len(e.PostFuncs); i++ {
		e.PostFuncs[i](service)
	}
}

func (e *ExtProps) doPostStartFuncs(ctx actor.Context) {
	if e.PostFuncs == nil {
		return
	}
	for i := 0; i < len(e.PostStartFuncs); i++ {
		e.PostStartFuncs[i](ctx)
	}
}

// NewServicePropsWithNewScheDisp 创建一个独立service
func NewServicePropsWithNewScheDisp(producer actor.Producer, name string) (*actor.Props, *ExtProps) {
	ext := &ExtProps{}

	disp := dispp.NewScheDisp(name)
	disp.Start()

	// can continue config like props.Configure(actor.WithDispatcher())
	props := actor.PropsFromProducer(func() actor.Actor {
		a := producer()
		s, ok := a.(IService)
		if !ok {
			l.Log.Warnf("必须实现 IService: %v", s)
			return nil
		}

		s.SetExtProps(ext)
		s.SetRunService(disp.GetRunService())

		if ext.PostFuncs != nil {
			ext.doPostFuncs(s)
		}

		if ext.APIs != "" {
			initDispatcherWithAPIs(s, ext.APIs)
		}

		if ext.Dispatcher != nil {
			s.SetAPIDispatcher(ext.Dispatcher)
		}

		s.OnCreate()
		return a
	}, actor.WithDispatcher(disp), actor.WithMailbox(mailbox.Producer(20)))

	return props, ext
}

func initDispatcherWithAPIs(s IService, apis string) {
	subs := strings.Split(apis, ",")
	if len(subs) == 0 {
		return
	}

	cols := make([]*apientry.APICollection, 0)
	for _, v := range subs {
		cols = append(cols, registry.Registry.GetCollection(v))
	}

	d := NewDispatcher(cols...)
	s.SetAPIDispatcher(d)
}
