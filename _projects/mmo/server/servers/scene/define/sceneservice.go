package define

import (
	"github.com/dfklegend/cell2/utils/event/light"
)

type ISceneService interface {
	GetEvents() *light.EventCenter
}
