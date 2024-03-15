package logic_old

import (
	"fmt"

	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/waterfall"

	mymsg "mmo/messages"
	"mmo/messages/cproto"
	"mmo/servers/logic_old/fightmode"
)

/*
	StartFight进入场景
	OnLeaveScene离开场景

	一次只能进入一个场景
	玩家战斗中，离线处理
	安全考虑应该不准离线，支持断线重连(TODO)

	IsInScene 防卡死，极端情况下，会去查询是否存在


	以后扩展不同的战斗模式，可以通过fightMode抽象来完成
	// 可以向客户端发送额外的消息
	// 可以向场景服推送额外的初始化战斗消息
*/

func (p *Player) StartFight(mode fightmode.ISceneFightMode, ns *service.NodeService, cb func(succ bool)) {
	// . 检查当前状态是否允许战斗
	// . 请求分配战斗服
	// . 请求进入战斗服

	// 后续cfgId可以知道是什么战斗类型
	cfgId := 1
	var allocAck *mymsg.SMAllocSceneAck
	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			if p.IsInScene() {
				p.ReqLeaveCurScene(func(succ bool) {
					callback(!succ)
				})
				return
			}
			callback(false)
		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {

			app.Request(ns, "scenem.remote.allocscene", nil, &mymsg.SMAllocScene{
				UId:   p.UId,
				CfgId: int32(cfgId),
			}, func(err error, raw interface{}) {
				if err != nil {
					l.L.Errorf("scenem allocscene failed")
					callback(true)
					return
				}

				ack := raw.(*mymsg.SMAllocSceneAck)
				if ack == nil {
					l.L.Errorf("bad SMAllocSceneAck")
					callback(true)
				}

				allocAck = ack

				mode.InitSceneInfo(ack.ServiceId, ack.SceneId, ack.Token)
				callback(false)
			})

		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			// TODO: 给场景推送战斗初始化信息
			if mode != nil {
				mode.SendInitDataToScene(func() {
					callback(false)
				})
				return
			}
			callback(false)
		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			// TODO: 给角色推送战斗额外信息
			if mode != nil {
				mode.SendInitDataToClient(func() {
					callback(false)
				})
				return
			}
			callback(false)
		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {

			l.L.Infof("SMAllocSceneAck: %v", allocAck)
			app.Request(ns, "scene.remote.enter", allocAck.ServiceId, &mymsg.SceneEnter{
				UId:     p.UId,
				SceneId: allocAck.SceneId,
				Token:   allocAck.Token,
				FrontId: p.FrontId,
				NetId:   p.NetId,
				LogicId: ns.Name,
				CfgId:   int32(cfgId),
			}, func(err error, raw any) {
				if err != nil {
					l.L.Errorf("scene enter error: %v", err)
					callback(true)
					return
				}

				p.onEnterScene(allocAck.SceneId, allocAck.ServiceId, int32(cfgId))
				l.L.Infof("scene enter succ")
				callback(false)
			})

		}).
		Final(func(err bool, args ...interface{}) {
			cb(!err)
		}).
		Do()
}

func (p *Player) OnGotFightReward(money int, exp int) {
	p.AddReward(money, exp)
	p.PushBattleLog(fmt.Sprintf("战斗胜利，获取: %v经验 %v金币", exp, money))
	p.PushCharInfo()
}

func (p *Player) PushBattleLog(str string) {
	app.PushMessageById(p.ns, p.FrontId, p.NetId, "battlelog", &cproto.BattleLog{
		Log: str,
	})
}

func (p *Player) SetCurFight(mode fightmode.ISceneFightMode) {
	p.curFight = mode
}

func (p *Player) GetCurFight() fightmode.ISceneFightMode {
	return p.curFight
}
