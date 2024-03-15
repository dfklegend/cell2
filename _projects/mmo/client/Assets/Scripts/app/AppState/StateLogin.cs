using Phoenix.Utils;

namespace Phoenix.Game
{  
    [IntType((int)eAppState.Login)]
    public class AppStateLogin : BaseAppState
    {
        public override void OnEnter()
        {
            Log.LogCenter.Default.Debug("Login.OnEnter");            
            UIMgr.It.OpenPanel("PanelLogin").Show();
            bindEvents(true);
        }

        public override void OnLeave()
        {
            UIMgr.It.OpenPanel("PanelLogin").Destroy();
            bindEvents(false);
        }

        private void bindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;            
            events.Bind(Card.EventDefine.LoginSucc, onLoginSucc, bind);
        }

        public override void Update()
        {            
        }

        private void onLoginSucc(params object[] args)
        {
            var e = args[0] as Card.HEventLoginSucc;

            // ��Ҫ��ȡ����            
            var client = ClientApp.It.client;
            client.ReqCharInfo(() => {
                // �ȴ�loadscene֪ͨ
                client.OpenCamera();
                //ClientApp.It.stateCtrl.ChangeState((int)eAppState.CardMain);
            });            
        }
    }    
}

