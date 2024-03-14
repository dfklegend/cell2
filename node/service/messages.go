package service

import (
	"github.com/dfklegend/cell2/node/config"
)

type StartServiceCmd struct {
	Name string
	Info *config.ServiceInfo
	Args []interface{}
}
