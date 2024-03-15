package script

var (
	fnLuaProviderCreator ProviderCreator
	fnGoProviderCreator  ProviderCreator
)

func SetLuaProviderCreator(fn ProviderCreator) {
	fnLuaProviderCreator = fn
}

// CreateLuaProvider
// 	*env.ScriptEnvData
func CreateLuaProvider(args ...any) IScriptProvider {
	return fnLuaProviderCreator(args...)
}

func SetGoProviderCreator(fn ProviderCreator) {
	fnGoProviderCreator = fn
}

func CreateGoProvider(args ...any) IScriptProvider {
	return fnGoProviderCreator(args...)
}
