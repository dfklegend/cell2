package center

import (
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/common"

	"mmo/common/define"
	mymsg "mmo/messages"
)

// 体验优化，登录时，等待老的链接下线，自动登录，而不是返回错误

type KickWaitTask struct {
	uid       int64
	frontId   string
	netId     uint32
	cbFunc    apientry.HandlerCBFunc
	startTime int64
}

func (t *KickWaitTask) Do(mgr *PlayerMgr) {
	mgr.ReqLogin(t.uid, t.frontId, t.netId, false, t.cbFunc)
}

func (t *KickWaitTask) Cancel(mgr *PlayerMgr) {
	// 取消掉上一个任务
	mgr.kickByNetId(t.uid, t.frontId, t.netId)
	//
	apientry.CheckInvokeCBFunc(t.cbFunc, nil, &mymsg.CenterReqLoginAck{
		Code: int32(define.ErrSystemBusy),
	})
}

func newKickWaitTask(uid int64, frontId string, netId uint32, cbFunc apientry.HandlerCBFunc) *KickWaitTask {
	return &KickWaitTask{
		uid:       uid,
		frontId:   frontId,
		netId:     netId,
		cbFunc:    cbFunc,
		startTime: common.NowMs(),
	}
}

type KickWaitTaskMgr struct {
	tasks            map[int64]*KickWaitTask
	ns               *service.NodeService
	mgr              *PlayerMgr
	nextCheckExpired int64
}

func newKickWaitTaskMgr(ns *service.NodeService, mgr *PlayerMgr) *KickWaitTaskMgr {
	return &KickWaitTaskMgr{
		tasks: make(map[int64]*KickWaitTask),
		ns:    ns,
		mgr:   mgr,
	}
}

func (m *KickWaitTaskMgr) AddTask(uid int64, frontId string, netId uint32, cbFunc apientry.HandlerCBFunc) {
	task := m.tasks[uid]
	if task != nil {
		m.ns.GetLogger().Infof("kick wait task exsit, overwrite it")

		// 取消掉上一个任务
		task.Cancel(m.mgr)
	}

	task = newKickWaitTask(uid, frontId, netId, cbFunc)
	m.tasks[uid] = task
}

func (m *KickWaitTaskMgr) removeTask(uid int64) {
	delete(m.tasks, uid)
}

func (m *KickWaitTaskMgr) OnSessionClose(uid int64) {
	m.tryRemoveExpired()

	task := m.tasks[uid]
	if task == nil {
		return
	}

	// 调用
	m.ns.GetLogger().Infof("do kickwait task: %v", uid)
	task.Do(m.mgr)
	m.removeTask(uid)
}

func (m *KickWaitTaskMgr) tryRemoveExpired() {
	now := common.NowMs()
	if now < m.nextCheckExpired {
		return
	}
	m.nextCheckExpired = now + 3*1000

	var toRemove []int64
	for k, v := range m.tasks {
		if now > v.startTime+30*1000 {
			toRemove = make([]int64, 0)
			toRemove = append(toRemove, k)
		}
	}

	if toRemove == nil {
		return
	}

	for _, uid := range toRemove {
		m.removeTask(uid)
	}
}
