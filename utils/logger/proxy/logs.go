package proxy

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logs = NewLogs()
)

func GetLogs() *Logs {
	return logs
}

type Logs struct {
	prefix string
	logDir string

	logs  map[string]*LogProxy
	mutex *sync.RWMutex
}

func NewLogs() *Logs {
	return &Logs{
		prefix: "node",
		logDir: "",
		logs:   make(map[string]*LogProxy),
		mutex:  &sync.RWMutex{},
	}
}

func (l *Logs) AddOrCreate(name string) *LogProxy {
	return l.AddLog(name)
}

func (l *Logs) AddLog(name string) *LogProxy {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.logs[name] != nil {
		return l.logs[name]
	}

	log := NewLog()
	log.Init(name, logrus.InfoLevel)
	if l.logDir != "" {
		log.EnableFileLog(l.prefix, l.logDir)
	}

	l.logs[name] = log
	return log
}

func (l *Logs) GetLog(name string) *LogProxy {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.logs[name]
}

func (l *Logs) EnableFileLog(prefix, logDir string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.prefix = prefix
	l.logDir = logDir
	for _, v := range l.logs {
		v.EnableFileLog(prefix, logDir)
	}
}
