package proxy

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

type simpleLogFormat struct {
	Name            string
	TimestampFormat string
}

func (s *simpleLogFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(entry.Time.Format(s.TimestampFormat))
	b.WriteString(fmt.Sprintf("[%s]", entry.Level.String()))
	b.WriteString(fmt.Sprintf("[%s]", s.Name))
	b.WriteString(fmt.Sprintf("%s", entry.Message))
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func NewSimpleLogFormat(name string) *simpleLogFormat {
	return &simpleLogFormat{
		Name:            name,
		TimestampFormat: "[2006-01-02 15:04:05.000]",
	}
}
