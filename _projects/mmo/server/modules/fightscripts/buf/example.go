package buf

import (
	"mmo/modules/fight/script"
)

func init() {
	mgr.Register("example", func() script.IBufScript {
		return &exampleScript{}
	})
}

type exampleScript struct {
}

func (s *exampleScript) Init(args ...any) {
}

func (s *exampleScript) OnStart(proxy script.IBufProxy) {
}

func (s *exampleScript) OnEnd(proxy script.IBufProxy) {
}

func (s *exampleScript) OnTriggle(proxy script.IBufProxy) {
}
