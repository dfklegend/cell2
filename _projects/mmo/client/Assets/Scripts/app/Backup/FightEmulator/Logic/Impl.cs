using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;



namespace Phoenix.Game.FightEmulator
{	    
    public class FightEmulatorImpl : ILogicImpl
    {   
        public void Start()
        {
            loadAllTables();
            initModelMgr();


            loadData();
            // 创建两个角色
            FightSimulateCtrl.It.Prepare();  
        }

        private void loadData()
        {
            UserPersistData.It.LoadFromFile();           
            
            FightSimulateCtrl.It.Load();            
        }

        private void initModelMgr()
        {
            UnitModelMgr.It.Init(
                UIMgr.It.GetPanel<PanelFightChars>().GetCharRoot());
        }

        public void Update()
        {
            // 驱动战斗
            UnitModelMgr.It.Update();
        }

        public void Destroy()
        {
            FightSimulateCtrl.It.Destroy();
            UnitModelMgr.It.Clear();
        }

        private void loadAllTables()
        {
            TestDataMgr.It.Load();
            ItemDataMgr.It.Load();
            EquipDataMgr.It.Load();
            SkillDataMgr.It.Load();
            MonsterDataMgr.It.Load();
        }
    }
} // namespace Phoenix
