package impls

import (
	"log"

	"github.com/dfklegend/cell2/utils/bridge"
	"github.com/dfklegend/cell2/utils/bridge/test/info"
)

type printer struct {
	info info.IInfo
}

func newPrinter(center *bridge.Center) *printer {
	return &printer{
		info: center.Get("info.Create").(info.IInfo),
	}
}

func (p *printer) Print(info string) {
	log.Println(info + p.info.GetInfo())
}
