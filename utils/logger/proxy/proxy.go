package proxy

import (
	"fmt"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/dfklegend/cell2/utils/logger/interfaces"
	logruswrapper "github.com/dfklegend/cell2/utils/logger/logrus"
)

type LogProxy struct {
	Name   string
	Log    interfaces.Logger
	logger *logrus.Logger
}

func NewLog() *LogProxy {
	return &LogProxy{}
}

//	TODO: 控制台和文件formatter保持一致
func (p *LogProxy) Init(name string, level logrus.Level) *LogProxy {
	p.Name = name

	plog := logrus.New()
	f := NewSimpleLogFormat(name)
	plog.SetFormatter(f)
	plog.Level = level

	log := plog.WithFields(logrus.Fields{
		"source": name,
	})

	p.Log, p.logger = logruswrapper.NewWithFieldLogger(log), plog
	return p
}

func (p *LogProxy) SetLogLevel(level logrus.Level) *LogProxy {
	p.logger.SetLevel(level)
	return p
}

func (p *LogProxy) SetFormatter(formatter logrus.Formatter) *LogProxy {
	p.logger.SetFormatter(formatter)
	return p
}

func (p *LogProxy) EnableFileLog(prefix, logDir string) *LogProxy {
	path := fmt.Sprintf("%v/%v.%v.%v.log", logDir, prefix, p.Name, "%Y%m%d%H%M")
	writer, _ := rotatelogs.New(
		path,
		rotatelogs.WithLinkName(path),
		// 7 天过期
		rotatelogs.WithMaxAge(time.Duration(7*24*3600)*time.Second),
		// 1 day rotation time
		rotatelogs.WithRotationTime(time.Duration(24*3600)*time.Second),
		// 最大文件
		rotatelogs.WithRotationSize(5*1024*1024),
	)

	p.logger.Hooks.Add(lfshook.NewHook(
		writer,
		&logrus.TextFormatter{
			DisableQuote: true,
		},
	))
	return p
}
