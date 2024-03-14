package pomelo

import (
	"github.com/dfklegend/cell2/node/client/impls"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/pomelonet/server/acceptor"
	"github.com/dfklegend/cell2/pomelonet/server/session"
)

// 	-------------
type WSComponent struct {
	*service.BaseComponent
	acceptor acceptor.Acceptor
	cfg      *session.SessionConfig
	sessions *impls.ClientSessions
}

func NewWSComponent(sessions *impls.ClientSessions) *WSComponent {
	return &WSComponent{
		BaseComponent: service.NewBaseComponent(),
		sessions:      sessions,
		cfg:           session.NewSessionConfig(nil),
	}
}

func (c *WSComponent) Start(address string) {
	c.cfg.Impl = NewSessionsImpl(c.GetNodeService().GetRunService().GetScheduler(), c.sessions)
	acceptor := acceptor.NewWSAcceptor(address)
	c.acceptor = acceptor

	StartAcceptor(acceptor, c.cfg)
}

func (c *WSComponent) OnAdd() {
}

func (c *WSComponent) OnRemove() {
}

// 	-------------
