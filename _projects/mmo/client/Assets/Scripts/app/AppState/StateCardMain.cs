using Phoenix.Utils;
using UnityEngine;

namespace Phoenix.Game
{  
    [IntType((int)eAppState.CardMain)]
    public class AppStateCardMain : BaseAppState
    {
        public override void OnEnter()
        {
            Log.LogCenter.Default.Debug("CardMain.OnEnter");
               

            UIMgr.It.OpenPanel("PanelShortcut").Hide();
            UIMgr.It.OpenPanel("PanelCards").Hide();
            UIMgr.It.OpenPanel("PanelCreateCard").Hide();
            UIMgr.It.OpenPanel("PanelCardFight").Hide();
            UIMgr.It.OpenPanel("PanelCardDetail").Hide();
            UIMgr.It.OpenPanel("PanelEquipTip").Hide();
            UIMgr.It.OpenPanel("PanelSkillTip").Hide();
            UIMgr.It.OpenPanel("PanelItemBag").Hide();
            UIMgr.It.OpenPanel("PanelBagSelect").Hide();
            UIMgr.It.OpenPanel("PanelCardMain").Show();
            bindEvents(true);
        }        

        private void bindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;
            events.Bind(Card.EventDefine.ConnectionBroken, onConnectionBroken, bind);            
        }

        public override void OnLeave()
        {
            UIMgr.It.GetPanel("PanelCardMain").Destroy();
            UIMgr.It.GetPanel("PanelShortcut").Destroy();
            UIMgr.It.GetPanel("PanelCards").Destroy();
            UIMgr.It.GetPanel("PanelCreateCard").Destroy();
            UIMgr.It.GetPanel("PanelCardFight").Destroy();
            UIMgr.It.GetPanel("PanelCardDetail").Destroy();
            UIMgr.It.GetPanel("PanelEquipTip").Destroy();
            UIMgr.It.GetPanel("PanelSkillTip").Destroy();
            UIMgr.It.GetPanel("PanelBagSelect").Destroy();
            bindEvents(false);
        }

        public override void Update()
        {            
        }        

        private void onConnectionBroken(params object[] args)
        {
            UIMgr.It.GetPanel<PanelDialog>().ShowInfo($"网络断开");
            ClientApp.It.stateCtrl.ChangeState((int)eAppState.Login);
        }
    }    
}

