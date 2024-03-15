package scenem

import (
	"errors"
	"strings"

	as "github.com/dfklegend/cell2/actorex/service"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/app"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/waterfall"

	"mmo/common/consts"
	"mmo/common/define"
	mymsg "mmo/messages"
)

func init() {
	registry.Registry.AddCollection("scenem.remote").
		Register(&Entry{}, apientry.WithName("remote"), apientry.WithNameFunc(strings.ToLower))
}

func Visit() {
}

type Entry struct {
	api.APIEntry
}

// AllocScene
// 分配一个场景
func (e *Entry) AllocScene(ctx *as.RemoteContext, msg *mymsg.SMAllocScene, cbFunc apientry.HandlerCBFunc) {
	// . 分配场景参数
	// . 要求场景服务 创建场景
	// . 返回创建成功
	s := ctx.ActorContext.Actor().(*Service)
	l.L.Infof("%v AllocScene enter:%v", s.Name, msg.UId)

	info := s.Mgr.AllocScene(msg.CfgId)

	ns := s.GetNodeService()
	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...interface{}) {

			app.Request(ns, "scene.remote.allocscene", info.ServiceId, &mymsg.SAllocScene{
				SceneId: info.SceneId,
				Token:   int32(info.Token),
				CfgId:   msg.CfgId,
			}, func(err error, raw interface{}) {
				if err != nil {
					l.L.Errorf("scene.remote.allocscene error: %v", err)
					callback(true)
					return
				}

				s.Mgr.OnSceneCreateSucc(info)
				callback(false)
			})

		}).
		Final(func(err bool, args ...interface{}) {
			if err {
				apientry.CheckInvokeCBFunc(cbFunc, errors.New("error allocscene"), nil)
			} else {
				apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.SMAllocSceneAck{
					ServiceId: info.ServiceId,
					SceneId:   info.SceneId,
					Token:     int32(info.Token),
				})
			}
		}).
		Do()
}

func (e *Entry) Refresh(ctx *as.RemoteContext, msg *mymsg.SMRefresh, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)

	s.Mgr.OnServiceRefresh(msg.ServiceId, int(msg.SceneNum))

	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(define.Succ),
	})
}

func (e *Entry) FreeScene(ctx *as.RemoteContext, msg *mymsg.SMOnSceneEnd) {
	s := ctx.ActorContext.Actor().(*Service)
	s.Mgr.OnSceneEnd(msg.SceneId)
}

func (e *Entry) QueryScenes(ctx *as.RemoteContext, msg *mymsg.SMQueryScenes, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)

	scenes := s.Mgr.QueryScenes(100)
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.SMAckScenes{
		Scenes: scenes,
	})
}

func (e *Entry) QueryScene(ctx *as.RemoteContext, msg *mymsg.SMQueryScene, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)

	scene := s.Mgr.GetScene(msg.SceneId)
	if scene == nil {
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrCanNotFindScene, nil)
		return
	}
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.SMQuerySceneAck{
		ServiceId: scene.ServiceId,
		SceneId:   scene.SceneId,
		CfgId:     scene.CfgId,
		Token:     int32(scene.Token),
	})
}

func (e *Entry) ReqSceneByCfgId(ctx *as.RemoteContext, msg *mymsg.SMReqSceneByCfgId, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)

	scene := s.Mgr.ReqSceneByCfgId(msg.CfgId)
	if scene == nil {
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrCanNotFindScene, nil)
		return
	}
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.SMQuerySceneAck{
		ServiceId: scene.ServiceId,
		SceneId:   scene.SceneId,
		CfgId:     scene.CfgId,
		Token:     int32(scene.Token),
	})
}
