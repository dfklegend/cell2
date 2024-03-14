package golua

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"errors"

	libs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

// compiledLuaProtoCache 编译好的lua proto缓存
var compiledLuaProtoCache map[string]*lua.FunctionProto
var userLuaPath string

func init() {
	InitLuaPathAndCompile("", false)
	compiledLuaProtoCache = make(map[string]*lua.FunctionProto)
}

// SetUserLuaPath 和编译功能分开比较合适，后续优化
// 如果一个工程目录下多个不同用途的lua脚本，用不同的入口函数来做区分
func setUserLuaPath(luaPath string) {
	if userLuaPath == luaPath {
		return
	}
	userLuaPath = luaPath
	switch userLuaPath {
	case ".":
		lua.LuaPathDefault = fmt.Sprintf("%s;./?.lua", lua.LuaPathDefault)
	case "":
		lua.LuaPathDefault = fmt.Sprintf("%s;?.lua", lua.LuaPathDefault)
	default:
		lua.LuaPathDefault = fmt.Sprintf("%s;%s/?.lua", lua.LuaPathDefault, luaPath)
	}
	os.Setenv(lua.LuaPath, lua.LuaPathDefault)
}

func InitLuaPathAndCompile(luaPath string, compile bool) {
	setUserLuaPath(luaPath)
	if compile {
		compileLuaFiles(luaPath)
	}
}

func CompileLuaFiles(luaPath string) {
	compileLuaFiles(luaPath)
}

type LuaEngine struct {
	L *lua.LState
}

func NewLuaEngine() *LuaEngine {

	engine := LuaEngine{
		L: lua.NewState(),
	}
	// 引入buildin库
	engine.L.OpenLibs()

	return &engine
}

// LoadModule 模块加载
func (e *LuaEngine) LoadModule(module string, loader func(L *lua.LState) int) {
	e.L.PreloadModule(module, loader)
}

// LoadGopherLuaLibs gopher-lua-libs
func (e *LuaEngine) LoadGopherLuaLibs() {
	libs.Preload(e.L)
}

func (e *LuaEngine) Close() {
	e.L.Close()
}

func (e *LuaEngine) DoLuaString(luaCode string) error {
	return e.L.DoString(luaCode)
}

func (e *LuaEngine) DoLuaFile(luaFile string) error {

	var fileName string
	if strings.HasSuffix(luaFile, ".lua") {
		fileName = luaFile
	} else {
		fileName = fmt.Sprintf("%s.lua", luaFile)
	}
	if IsCompiled(fileName) {
		return doCompiledLuaProto(e.L, compiledLuaProtoCache[fileName])
	}
	if "" == userLuaPath {
		luaFile = fileName
	} else {
		luaFile = fmt.Sprintf("%s/%s", userLuaPath, fileName)
	}
	return checkError(e.L.DoFile(luaFile))
}

// DoLuaMethod call a func in luaFile
func (e *LuaEngine) DoLuaMethod(luaFile string, method string, args ...any) error {
	if err := e.DoLuaFile(luaFile); err != nil {
		panic(err)
	}

	fn := e.L.Env.RawGetString(method)

	if fn == nil || fn.Type() != lua.LTFunction {
		panic(errors.New(fmt.Sprintf("not found lua method! %s:%s", luaFile, method)))
	}

	return Call(e.L, fn, args...)
	//var argsArr []lua.LValue
	//
	//for _, arg := range args {
	//	argsArr = append(argsArr, luar.New(e.L, arg))
	//}
	//
	//err := e.L.CallByParam(lua.P{
	//	Fn:      fn,
	//	NRet:    0,
	//	Protect: true,
	//	Handler: nil,
	//}, argsArr...)
	//
	//return checkError(err)
}

func (e *LuaEngine) DoLuaMethodWithResult(luaFile string, method string, args ...any) (lua.LValue, error) {

	if err := e.DoLuaFile(luaFile); err != nil {
		panic(err)
	}

	fn := e.L.Env.RawGetString(method)

	if fn == nil || fn.Type() != lua.LTFunction {
		panic(errors.New(fmt.Sprintf("not found lua method! %s:%s", luaFile, method)))
	}

	return CallWithResult(e.L, fn, args...)
	//var argsArr []lua.LValue
	//
	//for _, arg := range args {
	//	argsArr = append(argsArr, luar.New(e.L, arg))
	//}
	//
	//err := e.L.CallByParam(lua.P{
	//	Fn:      fn,
	//	NRet:    1,
	//	Protect: true,
	//	Handler: nil,
	//}, argsArr...)
	//
	//if err != nil {
	//	checkError(err)
	//}
	//
	//ret := e.L.Get(-1)
	//e.L.Pop(1)
	//
	//return ret, err
}

//compileLuaFiles 编译lua脚本并缓存map[luaName,luaProto]
func compileLuaFiles(luaDir string) {
	fileInfos, err := ioutil.ReadDir(luaDir)
	if err != nil {
		panic(err)
	}
	for _, fileInfo := range fileInfos {

		fileName := fileInfo.Name()

		if fileInfo.IsDir() {
			compileLuaFiles(fmt.Sprintf("%s/%s", luaDir, fileName))
			continue
		}

		if !strings.HasSuffix(fileName, ".lua") {
			continue
		}

		luaFile := fmt.Sprintf("%s/%s", luaDir, fileName)

		proto, err := compileLuaFile(luaFile)
		if err != nil {
			panic(err)
		}

		// cache lua compiled proto
		fileName = strings.ReplaceAll(luaFile, fmt.Sprintf("%s/", userLuaPath), "")
		compiledLuaProtoCache[fileName] = proto
	}
}

//IsCompiled 检查lua代码缓存是否存在
func IsCompiled(luaName string) bool {
	return compiledLuaProtoCache[luaName] != nil
}

func compileLuaFile(luaFile string) (*lua.FunctionProto, error) {
	file, err := os.Open(luaFile)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, luaFile)
	if err != nil {
		return nil, checkError(&lua.ApiError{Type: lua.ApiErrorSyntax, Object: lua.LString(err.Error()), Cause: err})
	}
	proto, err := lua.Compile(chunk, luaFile)
	if err != nil {
		return nil, checkError(err)
	}
	return proto, nil
}

func doCompiledLuaProto(L *lua.LState, proto *lua.FunctionProto) error {
	lfunc := L.NewFunctionFromProto(proto)
	L.Push(lfunc)
	return checkError(L.PCall(0, lua.MultRet, nil))
}

func stackTrace(apiError *lua.ApiError, level int) string {
	//got "exception/example.lua line:13(column:3) near 'end':   syntax error"
	line := apiError.Object.String()
	//need "exception/example.lua:12"
	file := strings.Split(line, " line:")[0]
	lineNum := substringBetween(line, "line:", "(column")
	n, _ := strconv.Atoi(lineNum)
	var buf []string
	header := "stack traceback:"
	buf = append(buf, fmt.Sprintf("\t%v", fmt.Sprintf("%s:%d", file, n-1))) //行号要减去1
	buf = append(buf, fmt.Sprintf("\t%v: %v", "[G]", "?"))
	buf = buf[intMax(0, intMin(level, len(buf))):len(buf)]
	if len(buf) > 20 {
		newbuf := make([]string, 0, 20)
		newbuf = append(newbuf, buf[0:7]...)
		newbuf = append(newbuf, "\t...")
		newbuf = append(newbuf, buf[len(buf)-7:len(buf)]...)
		buf = newbuf
	}
	return fmt.Sprintf("%s\n%s", header, strings.Join(buf, "\n"))
}

func intMin(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

//got "exception/example.lua line:13(column:3) near 'end':   syntax error"
//need "exception/example.lua:12"
func substringBetween(source string, left string, right string) string {
	return strings.Split(strings.Split(source, left)[1], right)[0]
}

func checkError(err error) error {
	if err == nil {
		return nil
	}

	var apiError *lua.ApiError
	if errors.As(err, &apiError) {
		if apiError.Type == lua.ApiErrorSyntax {
			apiError.StackTrace = stackTrace(apiError, 0)
		}
		panic(apiError)
	}
	return err
}
