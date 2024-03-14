package impls

import (
	"github.com/dfklegend/cell2/utils/bridge"
)

func Register(c *bridge.Center) {
	c.Register("printer.Create", func(args ...any) any {
		if len(args) == 0 {
			return nil
		}
		center := args[0].(*bridge.Center)
		return newPrinter(center)
	})
}
