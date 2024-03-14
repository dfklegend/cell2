package impls

import (
	"github.com/dfklegend/cell2/node/builtin/channel"
	"github.com/dfklegend/cell2/node/service"
)

// 	-------------
type ChannelComponent struct {
	*service.BaseComponent
	cs *channel.Service
}

func NewChannelComponent() *ChannelComponent {
	return &ChannelComponent{
		BaseComponent: service.NewBaseComponent(),
	}
}

func (c *ChannelComponent) GetCS() *channel.Service {
	return c.cs
}

func (c *ChannelComponent) OnAdd() {
	c.cs = channel.NewChannelService(c.GetNodeService())
}

func (c *ChannelComponent) OnRemove() {
}

// 	-------------
