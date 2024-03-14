package apientry

import "github.com/dfklegend/cell2/utils/serialize"

// 提供定义配置的方法
type (
	options struct {
		groupName    string              // component groupName
		nameFunc     func(string) string // rename group name and handler func name
		schedName    string              // schedName groupName
		serializer   serialize.Serializer
		serializeRet bool
	}

	// Option used to customize handler
	Option func(options *options)
)

// WithGroupName used to rename component groupName
func WithGroupName(name string) Option {
	return func(opt *options) {
		opt.groupName = name
	}
}

func WithInnerGroupName() Option {
	return func(opt *options) {
		opt.groupName = InnerGroupName
	}
}

// WithName 保持兼容
func WithName(name string) Option {
	return func(opt *options) {
		opt.groupName = name
	}
}

// WithNameFunc override handler groupName by specific function
// such as: strings.ToUpper/strings.ToLower
func WithNameFunc(fn func(string) string) Option {
	return func(opt *options) {
		opt.nameFunc = fn
	}
}

// WithSchedulerName set the groupName of the service scheduler
func WithSchedulerName(name string) Option {
	return func(opt *options) {
		opt.schedName = name
	}
}

func WithSerializer(serializer serialize.Serializer) Option {
	return func(opt *options) {
		opt.serializer = serializer
	}
}

func WithSerializeRet(do bool) Option {
	return func(opt *options) {
		opt.serializeRet = do
	}
}
