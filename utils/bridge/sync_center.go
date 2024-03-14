package bridge

import (
	"sync"
)

// SyncCenter
// routine safe
type SyncCenter struct {
	fns   map[string]FnGet
	mutex sync.Mutex
}

func NewSyncCenter() *SyncCenter {
	return &SyncCenter{
		fns: map[string]FnGet{},
	}
}

func (c *SyncCenter) Register(name string, fn FnGet) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.fns[name] = fn
}

func (c *SyncCenter) Get(name string, args ...any) any {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	fn := c.fns[name]
	if fn == nil {
		return nil
	}
	return fn(args...)
}
