using Phoenix.Core;
using UnityEngine;
using Phoenix.Entity;

namespace Phoenix.Game
{    
    public class AppModule : IAppComponent
    {   
        public int GetPriority()
        {
            return AppComponentPriority.AboveNormal + 1;
        }
        
        public void Start()
        {
            Platform.PlatformUtils.InitPlatform();
            Application.targetFrameRate = 60;
            SystemsUtils.RegisterAll(Systems.It);
            // init ui
            UIMgr.It.Init();            

            var worlds = GameObject.Find("Root/Worlds");
            WorldMgr.It.Init(worlds?.transform);
            WorldMgr.It.CreateWorld();

            GameLogicCtrl.It.Init();
            Config.It.Init();
        }        

        public bool IsReady()
        {
            return true;
        }

        public void Stop()
        {

        }

        public bool IsStopped()
        {
            return false;
        }

        public void PrepareUpdate()
        {
        }

        public void StopUpdate()
        {
        }

        public void Update()
        {            
            WorldMgr.It.Update();
        }
    }
}

