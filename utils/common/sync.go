package common

type IMutex interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

// FakeMutex 提供空的锁，便于测试同步(替代掉sync.Mutex,sync.RWMutex)
type FakeMutex struct {
}

func (f *FakeMutex) Lock()    {}
func (f *FakeMutex) Unlock()  {}
func (f *FakeMutex) RLock()   {}
func (f *FakeMutex) RUnlock() {}
