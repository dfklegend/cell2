package telnetcmd

import (
	"fmt"
	"strings"

	"github.com/reiver/go-telnet"

	"github.com/dfklegend/cell2/node/app"
)

type session struct {
	handler *Handler
	running bool
	close   chan int
	Output  chan string
}

func newSession() *session {
	return &session{
		running: true,
		close:   make(chan int),
		Output:  make(chan string, 9),
	}
}

func (s *session) handle(handler *Handler, ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	s.handler = handler
	s.pushOutput(fmt.Sprintf("welcome to master of cluster: %v\n", app.Node.GetClusterCfg().Name))
	go s.loopWrite(w)

	// hander需要阻塞，telnet的具体连接的处理
	// 所以直接用loopRead循环
	s.loopRead(r)
}

func (s *session) loopRead(r telnet.Reader) {
	buf := make([]byte, 1)
	bytes := make([]byte, 0)
	for s.running {

		_, err := r.Read(buf)
		if err != nil {
			//l.L.Errorf("telnet read error: %v", err)
			s.stop()
			break
		}
		bytes = append(bytes, buf[0])
		str := string(bytes)
		//l.L.Infof("got :%v %v", buf[0], str)

		// 一行一行解析
		index := strings.Index(str, "\n")
		if index >= 0 {
			line := str[:index]
			bytes = bytes[len([]byte(line))+1:]
			//l.L.Infof("got :%v ", line)

			line = line[:len(line)-1]

			s.handler.pushInput(s, line)
		}
	}
}

func (s *session) loopWrite(w telnet.Writer) {
	// 看看有没数据输出
	for s.running {
		select {
		case output := <-s.Output:
			_, err := w.Write([]byte(output))
			if err != nil {
				//l.L.Errorf("telnet write error: %v", err)
				break
			}
		case <-s.close:
			break
		}
	}

}

func (s *session) pushOutput(str string) {
	if !s.running {
		return
	}

	// 加一个换行
	if strings.Index(str, "\n") < 0 {
		str += "\n"
	}
	// 替换为\r\n
	str = strings.ReplaceAll(str, "\n", "\r\n")

	s.Output <- str
}

func (s *session) stop() {
	s.running = false
	close(s.close)
}
