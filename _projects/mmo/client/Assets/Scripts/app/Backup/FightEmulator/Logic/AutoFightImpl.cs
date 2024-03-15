using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;



namespace Phoenix.Game.FightEmulator
{	    
    // 创建一个角色，在场景内自由移动，自由打怪
    public class AutoFightImpl : ILogicImpl
    {
        int _spawned = 0;
        TimerID _timerSpawn;

        public void Start()
        {
            loadAllTables();
            initModelMgr();
            
            FightCtrl.It.Prepare();
            DisplayWorld.It.Init();

            //BuilderUtil.CreateWarrior(FightCtrl.It.GetWorld(), 100, 0);

            //int idPlayer = 1;
            //DisplayWorld.It.SetFocusEntity(idPlayer);
            //FightCtrl.It.PreparePlayer(idPlayer);

            //_timerSpawn = AppEnv.GetRunEnv().timer.AddTimer(OnTimerSpawn, 1f, 1f);

            bindEvents(true);
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
        }

        public void Destroy()
        {
            bindEvents(false);
            AppEnv.GetRunEnv().timer.Cancel(_timerSpawn);


            _timerSpawn = null;
            DisplayWorld.It.Destroy();
            FightCtrl.It.Destroy();
            UnitModelMgr.It.Clear();
        }

        private void bindEvents(bool bind)
        {
            var events = GlobalEvents.It.events;

            events.Bind(Card.EventDefine.TestSnapshot, onTestSnapshot, bind);
            events.Bind(Card.EventDefine.MoveTo, onMoveTo, bind);
            events.Bind(Card.EventDefine.UnitLeave, onUnitLeave, bind);
            events.Bind(Card.EventDefine.Attack, onAttack, bind);
            events.Bind(Card.EventDefine.UnitRelive, onUnitRelive, bind);
            events.Bind(Card.EventDefine.ServerStartSkill, onServerStartSkill, bind);
            events.Bind(Card.EventDefine.ServerSkillHit, onServerSkillHit, bind);
        }

        private void loadAllTables()
        {
            TestDataMgr.It.Load();
            ItemDataMgr.It.Load();
            EquipDataMgr.It.Load();
            SkillDataMgr.It.Load();
            MonsterDataMgr.It.Load();
        }

        private void OnTimerSpawn(object[] ps)
        {            
            if (_spawned++ > 5)
                return;
            string[] enemies = { "狗头人lv5", "熊lv5" };

            var index = MathUtil.RandomI(0, enemies.Length - 1);
            BuilderUtil.CreateMonster(FightCtrl.It.GetWorld(),
                enemies[index], 95, 1, MathUtil.RandomF(-10, 10), MathUtil.RandomF(-10,10));            
        }

        private void onTestSnapshot(params object[] args)
        {
            var e = args[0] as Card.HEventTestSnapshot;

            int entityId = BuilderUtil.CreateMonster(FightCtrl.It.GetWorld(),
                "熊lv5", 95, 1, e.snapshot.Pos.X, e.snapshot.Pos.Z );
            ServerCtrl.It.Add(e.snapshot.Id, entityId);

            var c = FightCtrl.It.GetChar(entityId);
            if (c == null)
                return;
            c.SetName(e.snapshot.Name);
            var attrs = c.attrs;
            attrs.GetAttr(AttrDefine.HPMax).Base.baseValue = e.snapshot.HPMax;
            attrs.GetAttr(AttrDefine.HP).Base.baseValue = e.snapshot.HP;
            CardUtils.SyncInfo(c);

            Log.LogCenter.Default.Debug("{0} {1}/{2}", e.snapshot.Id, e.snapshot.HP, e.snapshot.HPMax);
        }

        private void onMoveTo(params object[] args)
        {
            var e = args[0] as Card.HEventMoveTo;

            var entityId = ServerCtrl.It.GetEntityId(e.msg.Id);
            var c = FightCtrl.It.GetChar(entityId);
            if (c == null)
                return;
            c.moveCtrl.MoveTo(e.msg.Tar.X, e.msg.Tar.Z);
        }

        private void onUnitLeave(params object[] args)
        {
            var e = args[0] as Card.HEventUnitLeave;

            var entityId = ServerCtrl.It.GetEntityId(e.msg.Id);
            FightCtrl.It.DestroyEntity(entityId);            
        }

        private void onAttack(params object[] args)
        {
            // 转化为FightImpl的attack事件
            //var e = args[0] as Card.HEventAttack;
            //var src = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Id));
            //var tar = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Tar));
            //if (src == null || tar == null)
            //    return;

            //CardUtils.PlayAttack(src, tar);
            //CardUtils.PlaySimpleHit(src, tar, e.msg.Dmg, false);
            
            //CardUtils.CharSetHP(tar, e.msg.HPTar);                        
            //CardUtils.SyncHPBar(tar);
        }

        private void onUnitRelive(params object[] args)
        {
            var e = args[0] as Card.HEventRelive;
            var src = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Id));
            if (src == null)
                return;
            CardUtils.CharSetHP(src, e.msg.HP);
            CardUtils.SyncHPBar(src);
        }

        private void onServerStartSkill(params object[] args)
        {
            var e = args[0] as Card.HEventServerStartSkill;
            var src = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Id));
            var tar = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Tar));
            if (src == null || tar == null)
                return;
            CardUtils.PlayAttack(src, tar);
        }

        private void onServerSkillHit(params object[] args)
        {
            var e = args[0] as Card.HEventServerSkillHit;
            var src = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Id));
            var tar = FightCtrl.It.GetChar(ServerCtrl.It.GetEntityId(e.msg.Tar));
            if (src == null || tar == null)
                return;            
            CardUtils.PlaySimpleHit(src, tar, e.msg.Dmg, e.msg.Critical);

            CardUtils.CharSetHP(tar, e.msg.HPTar);
            CardUtils.SyncHPBar(tar);
        }
    }
} // namespace Phoenix
