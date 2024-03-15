package ctrl

func init() {

}

// AIOptions
// 可定义一些参数
type AIOptions struct {
	SearchInterval int64
	NextCanSearch  int64
}

func NewDefaultOptions() *AIOptions {
	return &AIOptions{
		SearchInterval: 0,
	}
}
