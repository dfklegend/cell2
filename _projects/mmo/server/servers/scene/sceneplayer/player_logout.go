package sceneplayer

import (
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/waterfall"

	mymsg "mmo/messages"
	"mmo/servers/scene/define"
)

/*
可以抽象一个重要异步事务的概念，比如支付等待结果
如果有等待中的重要异步事务，则不允许登出，切线等行为
那么在相应行为前，可以拒绝掉，如果有重要事务执行中，就不允许
比如
	匹配战斗前，切换场景前(远程场景)
*/

func (p *ScenePlayer) canLogoutNow() bool {
	// 致命错误，需要尽快下线
	if p.fatalErr && p.GetState() == define.Normal {
		now := common.NowMs()
		if now-p.lastTimeReqLogout < 60*1000 {
			return false
		}
		return true
	}
	// 正常状态
	if p.IsOnline() || p.GetState() != define.Normal {
		return false
	}
	now := common.NowMs()
	// 离线30s
	if now-p.timeOffline < 30*1000 || now-p.lastTimeReqLogout < 60*1000 {
		return false
	}

	// 其他条件，比如没有等待中的支付等
	return true
}

// 一般离线30s后，发起登出
func (p *ScenePlayer) tryLogout() {
	if !p.canLogoutNow() {
		return
	}
	p.beginLogout()
}

func (p *ScenePlayer) beginLogout() {
	// . 请求logout
	// . 存储player
	// . 通知onLogout
	// . 设置玩家可删除
	p.ChangeState(define.Logouting)
	p.lastTimeReqLogout = common.NowMs()

	ns := p.ns
	logger := ns.GetLogger()

	// 离开当前场景
	p.LeaveCurScene()

	logger.Infof("beginLogout: %v", p.uid)
	// 致命错误情况下，尽量保存下下线
	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...any) {
			if p.fatalErr {
				callback(false)
			} else {
				app.Request(ns, "logic.logicremote.reqlogout", p.logicId, &mymsg.ReqLogout{
					UId: p.uid,
				}, func(err error, raw interface{}) {
					if err != nil {
						logger.Errorf("logicremote.reqlogout failed: %v error: %v", p.uid, err)
						callback(true)
						return
					}
					callback(false)
				})
			}

		}).
		Next(func(callback waterfall.Callback, args ...any) {
			p.savePlayer(func(succ bool) {
				callback(!succ)
			})
		}).
		Next(func(callback waterfall.Callback, args ...any) {
			if p.fatalErr {
				callback(false)
			} else {
				logger.Infof("call logicremote.onlogout: %v", p.uid)
				app.Request(ns, "logic.logicremote.onlogout", p.logicId, &mymsg.OnLogout{
					UId: p.uid,
				}, func(err error, raw interface{}) {
					if err != nil {
						logger.Errorf("logicremote.onlogout failed: %v error: %v", p.uid, err)
						callback(true)
					}
					callback(false)
				})
			}
		}).
		Final(func(err bool, args ...any) {
			if err {
				// 等待下一次logout
				p.onLogoutFailed()
				return
			}
			p.ChangeState(define.WaitMgrDelete)
		}).
		Do()
}

func (p *ScenePlayer) onLogoutFailed() {
	logger := p.ns.GetLogger()
	logger.Errorf("player logout %v failed!", p.uid)

	if p.fatalErr {
		p.ChangeState(define.WaitMgrDelete)
		return
	}

	p.logoutErrTimes++
	if p.logoutErrTimes < 5 {
		// 多试几次
		p.ChangeState(define.Normal)
	} else {
		logger.Errorf("critical err, player %v cancel save and quit! ", p.uid)
		p.ChangeState(define.WaitMgrDelete)
	}
}

// SetFatalErr 致命错误，需要尽快logout
// 失去scenem服务的同步
// player向logic keepalive失败
func (p *ScenePlayer) SetFatalErr() {
	logger := p.ns.GetLogger()
	logger.Errorf("player %v SetFatalErr!", p.uid)

	p.fatalErr = true
	// 等待logout
}
