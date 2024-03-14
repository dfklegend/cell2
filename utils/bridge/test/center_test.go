package test

import (
	"testing"

	"github.com/dfklegend/cell2/utils/bridge"
	"github.com/dfklegend/cell2/utils/bridge/test/info/impls"
	printer2 "github.com/dfklegend/cell2/utils/bridge/test/printer"
	impls2 "github.com/dfklegend/cell2/utils/bridge/test/printer/impls"
)

func Test_Normal(t *testing.T) {
	center := bridge.NewCenter()

	impls.Register(center)
	impls2.Register(center)

	printer := center.Get("printer.Create", center).(printer2.IPrinter)
	printer.Print("some")
}
