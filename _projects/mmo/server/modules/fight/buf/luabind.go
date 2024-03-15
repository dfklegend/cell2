package buf

import (
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"

	"mmo/modules/fight/common"
)

// Proxy proxy for buf
type Proxy struct {
	buf *Buf
	// 用户数据
	// 脚本层使用
	UserData *lua.LTable
}

func newProxy(buf *Buf) *Proxy {
	return &Proxy{
		buf: buf,
	}
}

func (p *Proxy) GetId() common.BufId {
	return p.buf.id
}

func (p *Proxy) Owner() common.ICharProxy {
	return p.buf.owner.GetProxy()
}

func BindBuf(L *lua.LState) {
	L.SetGlobal("Buf", luar.NewType(L, Proxy{}))
}
