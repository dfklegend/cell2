package proxy

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_Normal(t *testing.T) {
	l := NewLog().Init("default", logrus.InfoLevel).
		EnableFileLog("test", "./logs")

	l.Log.Infof("hello!")
	l.Log.Errorf("error!")

	e := NewLog().Init("exception", logrus.InfoLevel).
		EnableFileLog("test", "./logs")
	e.Log.Errorf("error!")
}
