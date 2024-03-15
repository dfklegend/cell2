package gate

import (
	"github.com/dfklegend/cell2/node/app"
	cs "github.com/dfklegend/cell2/node/client/session"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/common"

	mymsg "mmo/messages"
)

func checkLogined(fs *cs.FrontSession) (uid int64, logined bool) {
	// uid设置了 = 登录完毕
	uid, ok := fs.Get("uid", 0).(int64)
	return uid, ok
}

// . 登录过程中，将在最后bind之前判断一下，如果掉线，通知center
// . 登录之后(通过Uid是否存在可知)，则直接通知center
func onSessionClose(ns *service.NodeService, fs *cs.FrontSession) {
	log := ns.GetLogger()
	log.Infof("routine:%v on session close: %v %v",
		common.GetRoutineID(), ns.Name, fs.GetNetId())

	// 如果在登录过程中，不发送(等到登录失败最后再发送)
	uid, logined := checkLogined(fs)
	if !logined {
		return
	}
	// 通知center
	sendSessionClose(ns, uid)
}

func sendSessionClose(ns *service.NodeService, uid int64) {
	log := ns.GetLogger()
	log.Infof("send session close: %v", uid)
	app.Request(ns, "center.centerremote.onsessionclose", nil, &mymsg.CenterOnSessionClose{
		UId: uid,
	}, func(err error, raw interface{}) {
		if err != nil {
			log.Errorf("centerremote.onsessionclose call failed: %v", err)
			return
		}
	})
}
