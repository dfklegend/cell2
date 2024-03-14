package impls

type info struct {
}

func newInfo() *info {
	return &info{}
}

func (c *info) GetInfo() string {
	return "hello from classa"
}
