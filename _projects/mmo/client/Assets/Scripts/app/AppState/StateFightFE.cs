using Phoenix.Utils;
using UnityEngine;
using UnityEngine.SceneManagement;
using Phoenix.Core;

namespace Phoenix.Game
{  
    [IntType((int)eAppState.FightFE)]
    public class AppStateFightFE : BaseAppState
    {
        public override void OnEnter()
        {
            Log.LogCenter.Default.Debug("Fight.OnEnter");

            UIMgr.It.GetPanel("PanelLoading").Show();
            Scene.SceneCtrl.It.LoadSceneAsync("fight", onFightSceneLoaded);
            bindEvents(true);            
        }

        private void bindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;
            events.Bind(Card.EventDefine.ConnectionBroken, onConnectionBroken, bind);
        }

        private void onFightSceneLoaded()
        {
            Log.LogCenter.Default.Debug("onFightSceneLoaded");
            UIMgr.It.OpenPanel("PanelFightFE").Show();
            UIMgr.It.OpenPanel("PanelFightChars").Show();
            UIMgr.It.OpenPanel("PanelMain").Hide();
            UIMgr.It.OpenPanel("PanelSelectEnemy").Hide();
            UIMgr.It.OpenPanel("PanelSelectRole").Hide();
            UIMgr.It.OpenPanel("PanelSelectEquip").Hide();
            UIMgr.It.OpenPanel("PanelItemBag").Hide();
            UIMgr.It.OpenPanel("PanelBagSelect").Hide();
            UIMgr.It.OpenPanel("PanelCharInfo").Hide();
            UIMgr.It.OpenPanel("PanelEquipTip").Hide();
            UIMgr.It.OpenPanel("PanelFightGraph").Hide();
            UIMgr.It.OpenPanel("PanelShortcut").Hide();
            UIMgr.It.OpenPanel("PanelEnterScene").Hide();
            UIMgr.It.OpenPanel("PanelDaySign").Hide();

            ItemPrefabFactory.It.Init();
            ItemListStyleFactory.It.Init();
            ItemIconStyleFactory.It.Init();

            WaitUIs waits = new WaitUIs();
            waits.Add("PanelFightFE");
            waits.Add("PanelFightChars");
            waits.Add("PanelSelectEnemy");
            waits.Add("PanelSelectRole");
            waits.Add("PanelSelectEquip");
            waits.Add("PanelItemBag");
            waits.Add("PanelCharInfo");
            waits.Add("PanelFightGraph");
            waits.Add("PanelEnterScene");
            waits.Add("PanelDaySign");
            waits.StartWait(() => 
            {
                Log.LogCenter.Default.Debug("AppStateFightFE waitUIs over");
                // UI都创建完毕
                GameLogicCtrl.It.ChangeImpl<Card.FightImpl>();
                Game.Card.InputSystem.It.Reset();

                var client = ClientApp.It.client;
                client.OnLoadSceneOver();
                Log.LogCenter.Default.Debug("OnLoadSceneOver");


                // 避免太快
                AppEnv.GetRunEnv().timer.AddTimer((args) => 
                {
                    UIMgr.It.GetPanel("PanelLoading").Hide();
                }, 0.2f, 0.0f);                
            });
        }        

        public override void OnLeave()
        {
            //SceneManager.UnloadSceneAsync("fight");
            GameLogicCtrl.It.DestroyImpl();
            Game.Card.InputSystem.It.OnPlayerLeave();

            UIMgr.It.OpenPanel("PanelFightFE").Hide();
            UIMgr.It.OpenPanel("PanelFightChars").Hide();
            UIMgr.It.OpenPanel("PanelFightGraph").Hide();

            bindEvents(false);
        }

        public override void Update()
        {
            Game.Card.InputSystem.It.Update();
        }

        private void onConnectionBroken(params object[] args)
        {
            UIMgr.It.GetPanel<PanelDialog>().ShowInfo($"网络断开");
            ClientApp.It.stateCtrl.ChangeState((int)eAppState.Login);
        }
    }    
}

