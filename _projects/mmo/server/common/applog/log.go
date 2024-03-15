package applog

import (
	"github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/logger/interfaces"
)

// 定义专门的log
var (
	DB       interfaces.Logger
	PlayerDB interfaces.Logger
)

func InitLogs() {
	DB = logger.NewLog("db")
	PlayerDB = logger.NewLog("playerdb")
}
