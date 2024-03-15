package redis

import (
	"errors"
	"strings"

	"github.com/dfklegend/cell2/actorex/service"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/serialize/proto/msgs"

	"mmo/common/define"
	mymsg "mmo/messages"
	"mmo/servers/db/dbop"
)

func init() {
	registry.Registry.AddCollection("db.redis").
		Register(&Entry{}, apientry.WithName("redis"), apientry.WithNameFunc(strings.ToLower))
}

type Entry struct {
	api.APIEntry
}

func (e *Entry) Init(ctx *service.RemoteContext, msg *mymsg.DBRedisInit, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	s.InitRedis(msg.Address)
	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
}

func (e *Entry) Auth(ctx *service.RemoteContext, msg *mymsg.DBAuth, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)

	s.TryCheckReady()
	uid := dbop.Auth(s.Client, msg.Username, msg.Password)

	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.DBAuthAck{
		UId: uid,
	})
}

func (e *Entry) LoadPlayer(ctx *service.RemoteContext, msg *mymsg.DBLoadPlayer, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)

	info, err := dbop.LoadPlayer(s.Client, msg.UId)
	if err != nil && dbop.IsNilError(err) {

		// 空值，新建角色
		l.L.Infof("redis, new player")
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.DBLoadPlayerAck{
			UId:       msg.UId,
			NewPlayer: true,
		})
		return
	}

	// redis失败，服务器连接不成功之类
	if err != nil {
		l.L.Errorf("redis error, LoadPlayer failed!")
		apientry.CheckInvokeCBFunc(cbFunc, errors.New("loadplayer failed"), &mymsg.DBLoadPlayerAck{})
		return
	}

	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.DBLoadPlayerAck{
		UId:  msg.UId,
		Info: info,
	})
}

func (e *Entry) SavePlayer(ctx *service.RemoteContext, msg *mymsg.DBSavePlayer, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)

	succ := dbop.SavePlayer(s.Client, msg.Info)

	code := define.Succ
	if !succ {
		code = define.ErrWithStr
	}
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(code),
	})
}

func (e *Entry) Notify(d *service.RemoteContext, msg *msgs.TestHello) {
}
