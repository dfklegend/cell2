package service

type BaseComponent struct {
	service INodeService
}

func NewBaseComponent() *BaseComponent {
	return &BaseComponent{}
}

func (c *BaseComponent) Init(service INodeService) {
	c.service = service
}

func (c *BaseComponent) GetNodeService() *NodeService {
	return c.service.(*NodeService)
}

func (c *BaseComponent) OnAdd() {
}

func (c *BaseComponent) OnRemove() {
}
