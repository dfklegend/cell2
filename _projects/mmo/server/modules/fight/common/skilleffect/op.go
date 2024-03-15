package skilleffect

// ----
// 效果类型
const (
	OpInvalid = iota
	OpAddBuf
	OpRemoveBuf
	OpAddBufWithLevel
	OpExample
)

var (
	OpNames = []string{
		"",
		"addbuf",
		"removebuf",
		"addbufwithlevel",
		"example",
	}
)
