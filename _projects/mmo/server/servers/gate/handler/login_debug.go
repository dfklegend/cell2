package gate

import (
	"time"

	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/waterfall"
)

/*
 * 登录问题主要在于客户端连接可能在任何时间断开，如果处理不仔细，可能造成角色多次登录情况，最终造成数据覆盖丢失
 * 如何测试登录异常，主要在登录流程中插入一个延时点，便于客户端断开连接，测试健壮性
 *
 * 简化处理
 * 登录步骤中，不用过多的判定是否掉线，而是最后一步在判定
 * 登录过程中的掉线通知(center)，推迟到最后一步再通知
 *
 * 测试
 * 登录过程中，关掉上一个连接，重新登录
 * 登录过程中，顶号
 */

const (
	loginTestDelayInterval = 0 // 增加的延时时间
)

func doTestDelay(ns *service.NodeService, callback waterfall.Callback) {
	// 测试延迟
	if loginTestDelayInterval == 0 {
		callback(false)
		return
	}

	log := ns.GetLogger()
	log.Infof("begin dtest delay")
	ns.GetRunService().GetTimerMgr().After(loginTestDelayInterval*time.Millisecond, func(args ...any) {
		log.Infof("dtest delay over")
		callback(false)
	})
}
