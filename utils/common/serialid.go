package common

import "sync/atomic"

/*
	提供简单自增id分配
	自动round
*/

type SerialIdService struct {
	nextId uint32
}

func (s *SerialIdService) AllocId() uint32 {
	// round
	v := atomic.AddUint32(&s.nextId, 1)
	if v == 0 {
		// 避免出现0
		return atomic.AddUint32(&s.nextId, 1)
	}
	return v
}

func (s *SerialIdService) doRoundAllocId() uint32 {
	return 0
}

func NewSerialIdService() *SerialIdService {
	return &SerialIdService{
		nextId: 1,
	}
}

type SerialIdService64 struct {
	nextId uint64
}

func (s *SerialIdService64) AllocId() uint64 {
	v := atomic.AddUint64(&s.nextId, 1)
	if v == 0 {
		// 避免出现0
		return atomic.AddUint64(&s.nextId, 1)
	}
	return v
}

func NewSerialIdService64() *SerialIdService64 {
	return &SerialIdService64{
		nextId: 1,
	}
}
