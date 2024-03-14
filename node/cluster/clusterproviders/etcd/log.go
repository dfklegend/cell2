package etcd

import (
	"github.com/dfklegend/cell2/utils/logger"
)

//var plog = log.New(log.DefaultLevel, "[CLUSTER] [ETCD]")
var plog = logger.Log

// SetLogLevel sets the log level for the logger
// SetLogLevel is safe to be called concurrently
