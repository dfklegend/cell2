package impls

import (
	"github.com/dfklegend/cell2/utils/bridge"
)

func Register(c *bridge.Center) {
	c.Register("info.Create", func(args ...any) any {
		return newInfo()
	})
}
