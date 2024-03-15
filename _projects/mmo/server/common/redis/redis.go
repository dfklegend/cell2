package redis

import (
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/go-redis/redis/v7"
)

//	redis component
type Component struct {
	*service.BaseComponent

	address string
	Client  *redis.Client
	IsReady bool
}

func NewComponent(address string) *Component {
	return &Component{
		BaseComponent: service.NewBaseComponent(),
		address:       address,
	}
}

func (c *Component) OnAdd() {
	//l.Log.Infof("handler on add")
	c.Client = redis.NewClient(&redis.Options{
		Addr:     c.address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := c.Client.Ping().Result()
	if err != nil {
		l.L.Errorf("start redis error: %v", err)
	} else {
		l.L.Infof("redis is ready: %v", c.address)
		c.IsReady = true
	}
}

func (c *Component) OnRemove() {
}
