package service

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

// 兼顾 如果使用nodeservice的log输出，自动带serviceid
type serviceLogFormat struct {
	logName         string
	timestampFormat string
}

func (s *serviceLogFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(entry.Time.Format(s.timestampFormat))
	b.WriteString(fmt.Sprintf("[%s]", entry.Level.String()))

	serviceName := entry.Data["service"]
	if serviceName != nil {
		b.WriteString(fmt.Sprintf("[%s]", serviceName))
	} else {
		b.WriteString(fmt.Sprintf("[%s]", s.logName))
	}

	b.WriteString(fmt.Sprintf("%s", entry.Message))
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func NewServiceLogFormat(logName string) *serviceLogFormat {
	return &serviceLogFormat{
		logName:         logName,
		timestampFormat: "[2006-01-02 15:04:05.000]",
	}
}
