using Phoenix.Utils;
using UnityEngine;
using UnityEngine.SceneManagement;

namespace Phoenix.Game
{  
    [IntType((int)eAppState.Fight)]
    public class AppStateFight : BaseAppState
    {
        public override void OnEnter()
        {
            Log.LogCenter.Default.Debug("Fight.OnEnter");

            Scene.SceneCtrl.It.LoadSceneAsync("fight", onFightSceneLoaded);            
        }

        private void onFightSceneLoaded()
        {
            Log.LogCenter.Default.Debug("onFightSceneLoaded");
            UIMgr.It.OpenPanel("PanelFight");

            WaitUIs waits = new WaitUIs();
            waits.Add("PanelFight");
            //waits.StartWait(() => { GameLogicCtrl.It.ChangeImpl<RandMove.RandMoveImpl>(); });
        }        

        public override void OnLeave()
        {
            SceneManager.UnloadSceneAsync("fight");
            GameLogicCtrl.It.DestroyImpl();
        }

        public override void Update()
        {            
        }
    }    
}

