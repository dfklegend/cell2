package utils

type TestTimeProvider struct {
	now int64
}

func NewTestTimeProvider() *TestTimeProvider {
	return &TestTimeProvider{}
}

func (p *TestTimeProvider) SetNow(now int64) {
	p.now = now
}

func (p *TestTimeProvider) NowMs() int64 {
	return p.now
}

func (p *TestTimeProvider) Update() {
}
