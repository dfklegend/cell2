package sceneplayer

import (
	"github.com/dfklegend/cell2/node/app"

	"mmo/common/define"
	mymsg "mmo/messages"
	define2 "mmo/servers/scene/define"
	"mmo/utils"
)

func (p *ScenePlayer) EnterPreNormal() {
	// 向logic二次确认，确认是否成功
	ns := p.ns
	logger := ns.GetLogger()
	app.Request(ns, "logic.logicremote.checkplayer2", p.logicId, &mymsg.CheckPlayer2{
		UId:     p.uid,
		SceneId: p.sceneId,
	}, func(err error, ret interface{}) {
		ack := utils.TryGetNormalAck(ret)
		if err != nil || (ack != nil && ack.Code != int32(define.Succ)) {
			logger.Errorf("checkplayer2 failed! player: %v", p.uid)
			if err != nil {
				logger.Errorf("checkplayer2 failed! player: %v, err: %v", p.uid, err)
			}
			p.ChangeState(define2.WaitMgrDelete)
			return
		}

		p.ChangeState(define2.Normal)
	})
}
