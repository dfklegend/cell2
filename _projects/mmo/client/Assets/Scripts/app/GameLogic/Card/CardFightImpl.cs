using UnityEngine;
using Phoenix.Core;

namespace Phoenix.Game.Card
{	    
    // 
    public class FightImpl : ILogicImpl
    {
        public void Start()
        {
            loadAllTables();
            initModelMgr();
            
            FightCtrl.It.Prepare();
            DisplayWorld.It.Init();
            
            bindEvents(true);
            DataCenter.It.fightStat.Start();
        }        

        private void initModelMgr()
        {
            UnitModelMgr.It.Init(
                UIMgr.It.GetPanel<PanelFightChars>().GetCharRoot());
        }

        public void Update()
        {
            // 驱动战斗
            DisplayWorld.It.Update();
            DataCenter.It.fightStat.Update();
        }

        public void Destroy()
        {
            DataCenter.It.fightStat.Stop();
            bindEvents(false);

       
            DisplayWorld.It.Destroy();
            FightCtrl.It.Destroy();
            UnitModelMgr.It.Clear();
            ServerCtrl.It.Clear();
        }

        private void bindEvents(bool bind)
        {
            var events = GlobalEvents.It.events;

            events.Bind(EventDefine.TestSnapshot, onTestSnapshot, bind);
            events.Bind(EventDefine.ExitSnapshot, onExitSnapshot, bind);
            events.Bind(EventDefine.MonsterSnapshot, onMonsterSnapshot, bind);
            events.Bind(EventDefine.MoveTo, onMoveTo, bind);
            events.Bind(EventDefine.StopMove, onStopMove, bind);
            events.Bind(EventDefine.UnitLeave, onUnitLeave, bind);
            events.Bind(EventDefine.Attack, onAttack, bind);
            events.Bind(EventDefine.UnitRelive, onUnitRelive, bind);
            events.Bind(EventDefine.ServerStartSkill, onServerStartSkill, bind);
            events.Bind(EventDefine.ServerSkillHit, onServerSkillHit, bind);
            events.Bind(EventDefine.ServerSkillBroken, onServerSkillBroken, bind);
            events.Bind(EventDefine.UnitAttrsChanged, onServerUnitAttrsChanged, bind);
        }

        private void loadAllTables()
        {
            
        }        

        private void onTestSnapshot(params object[] args)
        {
            var e = args[0] as HEventTestSnapshot;

            int entityId = BuilderUtil.CreatePlayer(FightCtrl.It.GetWorld(),
                95, 1, e.snapshot.Pos.X, e.snapshot.Pos.Z);
            ServerCtrl.It.Add(e.snapshot.Id, entityId);

            var c = FightCtrl.It.GetChar(entityId);
            if (c == null)
                return;
            c.SetName(e.snapshot.Name);
            c.SetHPMax(e.snapshot.HPMax);
            c.SetHP(e.snapshot.HP);
            c.SetSide(e.snapshot.Side);
            CardUtils.SyncInfo(c);
            CardUtils.SyncHPBar(c);            

            DataCenter.It.fightStat.AddUnitInfo(e.snapshot.Id, e.snapshot.Name);

            Log.LogCenter.Default.Debug("{0} {1}/{2}", e.snapshot.Id, e.snapshot.HP, e.snapshot.HPMax);
        }

        private void onExitSnapshot(params object[] args)
        {
            var e = args[0] as HEventExitSnapshot;

            ExitBuilder.It.CreateExit(FightCtrl.It.GetWorld(),
                new ExitCreateInfo() 
                {
                    pos = new Vector3(e.msg.Pos.X, 0, e.msg.Pos.Z),
                });
        }

        private void onMonsterSnapshot(params object[] args)
        {
            var e = args[0] as HEventMonsterSnapshot;

            int entityId = BuilderUtil.CreateMonster(FightCtrl.It.GetWorld(),
                e.msg.CfgId, 95, 1, e.msg.Pos.X, e.msg.Pos.Z);
            ServerCtrl.It.Add(e.msg.Id, entityId);

            var c = FightCtrl.It.GetChar(entityId);
            if (c == null)
                return;
            c.SetName(e.msg.Name);
            c.SetHPMax(e.msg.HPMax);
            c.SetHP(e.msg.HP);
            c.SetSide(e.msg.Side);
            CardUtils.SyncInfo(c);
            CardUtils.SyncHPBar(c);

            DataCenter.It.fightStat.AddUnitInfo(e.msg.Id, e.msg.Name);

            Log.LogCenter.Default.Debug("{0} {1}/{2}", e.msg.Id, e.msg.HP, e.msg.HPMax);
        }

        private void onMoveTo(params object[] args)
        {
            var e = args[0] as HEventMoveTo;

            var entityId = ServerCtrl.It.GetEntityId(e.msg.Id);
            var c = FightCtrl.It.GetChar(entityId);
            if (c == null)
                return;
            c.moveCtrl.MoveTo(e.msg.Tar.X, e.msg.Tar.Z);
        }

        private void onStopMove(params object[] args)
        {
            var e = args[0] as HEventStopMove;

            var entityId = ServerCtrl.It.GetEntityId(e.msg.Id);
            var c = FightCtrl.It.GetChar(entityId);
            if (c == null)
                return;
            // 可以看看误差范围
            c.moveCtrl.Stop();
        }

        private void onUnitLeave(params object[] args)
        {
            var e = args[0] as HEventUnitLeave;

            var entityId = ServerCtrl.It.GetEntityId(e.msg.Id);
            FightCtrl.It.DestroyEntity(entityId);

            DataCenter.It.fightStat.Stop();
        }

        private void onAttack(params object[] args)
        {            
        }

        private void onUnitRelive(params object[] args)
        {
            var e = args[0] as HEventRelive;
            var src = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Id));
            if (src == null)
                return;
            CardUtils.CharSetHP(src, e.msg.HP);
            CardUtils.SyncHPBar(src);
        }

        private void onServerStartSkill(params object[] args)
        {
            var e = args[0] as HEventServerStartSkill;
            var src = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Id));
            var tar = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Tar));
            if (src == null || tar == null)
                return;
            CardUtils.PlayAttack(e.msg.SkillId, src, tar);
        }

        private void onServerSkillHit(params object[] args)
        {
            var e = args[0] as HEventServerSkillHit;
            var src = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Id));
            var tar = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Tar));
            if (src == null || tar == null)
                return;
            CardUtils.PlaySimpleHit(e.msg.SkillId, src, tar, e.msg.Dmg, e.msg.Critical);
            DataCenter.It.fightStat.AddDmg(e.msg.Id, e.msg.Dmg);

            CardUtils.CharSetHP(tar, e.msg.HPTar);
            CardUtils.SyncHPBar(tar);
        }

        private void onServerSkillBroken(params object[] args)
        {
            var e = args[0] as HEventServerSkillBroken;
            var src = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Id));
            if (src == null)
                return;
            CardUtils.DoSkillBroken(src, e.msg.SkillId);
        }

        private void onServerUnitAttrsChanged(params object[] args)
        {
            var e = args[0] as HEventServerAttrsChanged;
            var src = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Id));
            if (src == null)
                return;
            var msg = e.msg;
            for (var i = 0; i < msg.Attrs.Count; i++)
            {
                var one = msg.Attrs[i];
                src.SetAttr(one.Index, one.Value);

                if(one.Index == (int)eAttrType.Energy)
                {
                    CardUtils.SyncMPBar(src);
                }
            }
        }
    }
} // namespace Phoenix
