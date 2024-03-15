package utils

import (
	"strconv"

	"github.com/dfklegend/cell2/node/client/impls"
	l "github.com/dfklegend/cell2/utils/logger"

	mymsg "mmo/messages"
)

func GetIdFromContext(ctx *impls.HandlerContext) int64 {
	sid := ctx.Session.GetID()
	uid, err := strconv.ParseInt(sid, 0, 64)
	if err != nil {
		l.L.Errorf("error uid: %v", sid)
		return 0
	}
	return uid
}

// TryGetNormalAck
// 有时候，对方返回是nil
func TryGetNormalAck(ret any) *mymsg.NormalAck {
	if ret != nil {
		return ret.(*mymsg.NormalAck)
	}
	return nil
}
