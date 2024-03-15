package bufdef

const (
	TarTypeInvalid = iota
	TarTypeOwner
	TarTypeOwnerTar
)

var (
	TarTypeNames = []string{
		"",
		"owner",
		"ownertar",
	}
)
