// Copyright (c) nano Author and TFG Co. All Rights Reserved.
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

package session

/*
 	FlowCtrl
	流量控制
	控制进入流量，如果进入流量持续的超过上限，会被强制关闭
 	如果超出了发送限制，暂缓发送

*/

type FlowCtrl struct {
	maxInSpeed  int32
	maxOutSpeed int32

	canSendSize int32
	canRecvSize int32

	// 上次恢复时间
	lastRestore int64
}

func NewFlowCtrl(maxIn int32, maxOut int32) *FlowCtrl {
	return &FlowCtrl{
		maxInSpeed:  maxIn,
		maxOutSpeed: maxOut,
	}
}

func (f *FlowCtrl) ReqSend(size int32) bool {
	return true
}

func (f *FlowCtrl) ReqProcess(size int32) bool {
	return true
}

func (f *FlowCtrl) tryRestore() {
}

// ReadCtrl
// 避免客户端上传流量太多，监测速度，报告状态
// 持续异常状态，关闭链接
type ReadCtrl struct {
}

// WriteCtrl
// 避免短时间写入太多，占用其他人带宽
type WriteCtrl struct {
	writeRecently int
}

func (w *WriteCtrl) ReqWrite(size int) bool {
	w.writeRecently += size
	return true
}
