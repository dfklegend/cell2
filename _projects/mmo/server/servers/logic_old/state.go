package logic_old

type ServiceState struct {
	retire *RetireMgr
}

func NewServiceState() *ServiceState {
	return &ServiceState{}
}

func (s *ServiceState) Init(retire *RetireMgr) {
	s.retire = retire
}

// CanWork 是否正常工作
func (s *ServiceState) CanWork() bool {
	// 非退休状态
	return s.retire.IsRunning()
}
