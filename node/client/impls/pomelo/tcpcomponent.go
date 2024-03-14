package pomelo

import (
	"github.com/dfklegend/cell2/node/client/impls"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/pomelonet/server/acceptor"
	"github.com/dfklegend/cell2/pomelonet/server/session"
)

// 	-------------
type TCPComponent struct {
	*service.BaseComponent
	acceptor acceptor.Acceptor
	cfg      *session.SessionConfig
	sessions *impls.ClientSessions
}

func NewTCPComponent(sessions *impls.ClientSessions) *TCPComponent {
	return &TCPComponent{
		BaseComponent: service.NewBaseComponent(),
		sessions:      sessions,
		cfg:           session.NewSessionConfig(nil),
	}
}

func (c *TCPComponent) Start(address string) {
	c.cfg.Impl = NewSessionsImpl(c.GetNodeService().GetRunService().GetScheduler(), c.sessions)
	acceptor := acceptor.NewTCPAcceptor(address)
	c.acceptor = acceptor

	StartAcceptor(acceptor, c.cfg)
}

func (c *TCPComponent) OnAdd() {
}

func (c *TCPComponent) OnRemove() {
}

// 	-------------
