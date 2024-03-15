using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using Phoenix.API;
using Phoenix.Log;
using Phoenix.Core;
using Phoenix.Game.Card;

namespace Network
{   
    // 推送的消息
    [APIService("client", "")]
    public class ServicePush : IAPIService
    {
        [APIFunc]
        public void kick(IContext context, Cproto.Kick msg)
        {
            LogCenter.Default.Debug($"被踢掉线, reason: {msg.Reason}");
        }

        [APIFunc]
        public void charinfo(IContext context, Cproto.CharInfo msg)
        {            
            LogCenter.Default.Debug($"got charinfo: {msg.ToString()}");
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventCharInfo(msg));
        }

        [APIFunc]
        public void enterbegin(IContext context, Cproto.EmptyArg msg)
        {
            LogCenter.Default.Debug($"enter begin");
            
        }

        [APIFunc]
        public void enterend(IContext context, Cproto.EmptyArg msg)
        {
            LogCenter.Default.Debug($"enter end");
        }

        [APIFunc]
        public void servercmd(IContext context, Cproto.ServerSystemCmd msg)
        {   
            //LogCenter.Default.Debug($"got servercmd: {msg.System}.{msg.Cmd}");
            Phoenix.Game.Systems.It.OnServerCmd(msg.System, msg.Cmd, msg.Args.ToByteArray());
        }

        [APIFunc]
        public void battlelog(IContext context, Cproto.BattleLog msg)
        {
            LogCenter.Default.Debug($"battle log: {msg.Log}");
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventBattleLog(msg.Log));
        }

        [APIFunc]
        public void testsnapshot(IContext context, Cproto.TestSnapshot msg)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventTestSnapshot(msg));
        }

        [APIFunc]
        public void exitsnapshot(IContext context, Cproto.ExitSnapshot msg)
        {
            //传送点
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventExitSnapshot(msg));
        }

        [APIFunc]
        public void monstersnapshot(IContext context, Cproto.MonsterSnapshot msg)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventMonsterSnapshot(msg));
        }

        [APIFunc]
        public void moveto(IContext context, Cproto.MoveTo msg)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventMoveTo(msg));
        }

        [APIFunc]
        public void stopmove(IContext context, Cproto.MoveTo msg)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventStopMove(msg));
        }

        [APIFunc]
        public void attack(IContext context, Cproto.Attack msg)
        {
            //HEventUtil.Dispatch(GlobalEvents.It.events, new Phoenix.Game.Card.HEventAttack(msg));
        }

        [APIFunc]
        public void startskill(IContext context, Cproto.StartSklill msg)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events, new Phoenix.Game.Card.HEventServerStartSkill(msg));
        }

        [APIFunc]
        public void skillhit(IContext context, Cproto.SkillHit msg)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventServerSkillHit(msg));
        }

        [APIFunc]
        public void skillbroken(IContext context, Cproto.SkillBroken msg)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventServerSkillBroken(msg));
        }

        [APIFunc]
        public void unitleave(IContext context, Cproto.UnitLeave msg)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventUnitLeave(msg));
        }

        [APIFunc]        
        public void loadscene(IContext context, Cproto.LoadScene msg)
        {
            // 服务器要求客户端载入场景
            DataCenter.It.SceneServer = msg.ServerId;
            DataCenter.It.SceneId = msg.SceneId;
            DataCenter.It.SceneCfgId = msg.CfgId;

            // 根据cfgId可以选择不同战斗模式
            ClientApp.It.stateCtrl.ChangeState((int)Phoenix.Game.eAppState.FightFE);            
        }

        [APIFunc]
        public void onchangescene(IContext context, Cproto.LoadScene msg)
        {
            // 切换了场景
            // 清除本地战斗对象
            // 载入场景
            if(ClientApp.It.stateCtrl.GetStateType() == (int)Phoenix.Game.eAppState.FightFE)
            {
                DataCenter.It.SceneServer = msg.ServerId;
                DataCenter.It.SceneId = msg.SceneId;
                DataCenter.It.SceneCfgId = msg.CfgId;

                // 目前无需载入场景                
                Phoenix.Game.GameLogicCtrl.It.DestroyImpl();
                Phoenix.Game.GameLogicCtrl.It.ChangeImpl<Phoenix.Game.Card.FightImpl>();

                var client = ClientApp.It.client;
                client.OnLoadSceneOver();

                HEventUtil.Dispatch(GlobalEvents.It.events, new HEventRefreshSceneInfo());
            }
        }

        [APIFunc]
        // 进入场景
        public void avatarenter(IContext context, Cproto.AvatarEnter msg)
        {
            DataCenter.It.SceneAvatarId = msg.Id;

            var idPlayer = ServerCtrl.It.GetEntityId(msg.Id);
            DisplayWorld.It.SetFocusEntity(idPlayer);
            FightCtrl.It.PreparePlayer(idPlayer);

            InputSystem.It.OnPlayerEnter(FightCtrl.It.player);
        }

        [APIFunc]        
        public void unitrelive(IContext context, Cproto.UnitRelive msg)
        {
            // 复活
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventRelive(msg));
        }

        [APIFunc]
        public void refreshcards(IContext context, Cproto.RefreshCards msg)
        {           
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventServerRefreshCards(msg));
        }

        [APIFunc]
        public void attrschanged(IContext context, Cproto.UnitAttrsChanged msg)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventServerAttrsChanged(msg));
        }
    }
}
