package baseapp

import (
	"github.com/dfklegend/cell2/baseapp/interfaces"
	l "github.com/dfklegend/cell2/utils/logger"
)

var (
	LaunchFactory = NewFacotry()
	defaultMode   interfaces.ILaunchMode
)

//	启动服务器

func LaunchAppWithMode(app interfaces.IApp, mode interfaces.ILaunchMode, finCB interfaces.FuncWithSucc) bool {
	mode.PrepareModules(app)
	app.Start(finCB)
	return true
}

func LaunchApp(app interfaces.IApp, modeName string, finCB interfaces.FuncWithSucc) bool {
	var mode interfaces.ILaunchMode
	if modeName != "" {
		mode = LaunchFactory.GetMode(modeName)
		if mode == nil {
			l.L.Warnf("can not find launchmode: %v fallback to default", modeName)
			mode = defaultMode
		}
	} else {
		mode = defaultMode
	}

	if mode == nil {
		l.L.Errorf("nil launchmode, no default too")
		return false
	}

	return LaunchAppWithMode(app, mode, finCB)
}

type FuncMode struct {
	f func(app interfaces.IApp)
}

func NewFuncMode(f func(app interfaces.IApp)) interfaces.ILaunchMode {
	return &FuncMode{
		f: f,
	}
}

func (f *FuncMode) PrepareModules(app interfaces.IApp) {
	f.f(app)
}

// 	工厂
type launchFactory struct {
	modes map[string]interfaces.ILaunchMode
}

func NewFacotry() *launchFactory {
	return &launchFactory{
		modes: make(map[string]interfaces.ILaunchMode),
	}
}

func (l *launchFactory) Register(name string, mode interfaces.ILaunchMode) {
	l.modes[name] = mode
}

func (l *launchFactory) GetMode(name string) interfaces.ILaunchMode {
	return l.modes[name]
}

func RegisterLaunchFunc(name string, f func(app interfaces.IApp)) {
	LaunchModeRegister(name, NewFuncMode(f))
}

func LaunchFunc(name string, f func(app interfaces.IApp)) {
	RegisterLaunchFunc(name, f)
}

func SetDefaultLaunchFunc(f func(app interfaces.IApp)) {
	defaultMode = NewFuncMode(f)
}

func LaunchModeRegister(name string, mode interfaces.ILaunchMode) {
	LaunchFactory.Register(name, mode)
}

func LaunchModeGet(name string) interfaces.ILaunchMode {
	return LaunchFactory.GetMode(name)
}
