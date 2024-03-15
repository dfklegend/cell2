using Phoenix.Utils;
using UnityEngine;

namespace Phoenix.Game
{  
    [IntType((int)eAppState.Init)]
    public class AppStateInit : BaseAppState
    {
        public override void OnEnter()
        {
            Log.LogCenter.Default.Debug("Init.OnEnter");

            Test();
            UIMgr.It.OpenPanel("PanelWelcome");
            UIMgr.It.OpenPanel("PanelFloatText");
            UIMgr.It.OpenPanel("PanelDialog").Hide();
            UIMgr.It.OpenPanel("PanelLoading").Hide();

            Card.CardUtils.loadAllPrefabs();
            Card.CardUtils.loadAllTables();
        }

        public override void OnLeave()
        {
        }

        public override void Update()
        {            
        }

        private void Test()
        {
            //Debug.LogAssertion("a assert");
            //throw new System.Exception("a exception");
        }
    }    
}

