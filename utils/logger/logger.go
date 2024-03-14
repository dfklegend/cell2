// Copyright (c) TFG Co. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package logger

// 方便接口
// 预定义了普通log和exception log
import (
	"github.com/sirupsen/logrus"

	"github.com/dfklegend/cell2/utils/logger/interfaces"
	"github.com/dfklegend/cell2/utils/logger/proxy"
)

// Log is the default logger
var (
	Log        interfaces.Logger = nil
	L          interfaces.Logger = nil
	defaultLog *proxy.LogProxy

	// 	Exception 异常log
	Exception interfaces.Logger = nil
	E         interfaces.Logger = nil
)

func init() {
	defaultLog = proxy.GetLogs().AddLog("default")
	Log = defaultLog.Log
	L = Log

	Exception = proxy.GetLogs().AddLog("exception").Log
	E = Exception
}

func EnableFileLog(prefix, logDir string) {
	proxy.GetLogs().EnableFileLog(prefix, logDir)
}

// SetLogger rewrites the default logger
func SetLogger(l interfaces.Logger) {
	if l != nil {
		Log = l
	}
}

func SetLogLevel(level logrus.Level) {
	defaultLog.SetLogLevel(level)
}

func SetDebugLevel() {
	SetLogLevel(logrus.DebugLevel)
}

func SetInfoLevel() {
	SetLogLevel(logrus.InfoLevel)
}

func SetWarnLevel() {
	SetLogLevel(logrus.WarnLevel)
}

func NewLog(name string) interfaces.Logger {
	return proxy.GetLogs().AddLog(name).Log
}

func GetLog(name string) interfaces.Logger {
	return proxy.GetLogs().GetLog(name).Log
}

func GetLogProxy(name string) *proxy.LogProxy {
	return proxy.GetLogs().GetLog(name)
}
