using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Log;

namespace Phoenix.Game
{
    [StringType("PanelLogin")]
    public class PanelLogin : BasePanel
    {
        private Button _btnLogin;
        private InputField _input;
        public override void OnReady()
        {
            SetDepth(100);
            base.OnReady();

            _btnLogin = _root.Find("BG/btnLogin").GetComponent<Button>();
            _btnLogin.onClick.AddListener(() =>
            {
                onLogin();
            });

            _input = _root.Find("BG/inputUser").GetComponent<InputField>();

            BindEvents(true);
        }

        private void onLogin()
        {
            // check input
            var user = _input.text;
            if (user == "")
            {
                LogCenter.Default.Debug($"username is empty");
                UIMgr.It.GetPanel<PanelDialog>().ShowInfo($"用户名为空");
                return;
            }

            var client = ClientApp.It.client;
            client.username = user;

            var ip = Config.It.main.GetString("main", "server", "127.0.0.1");
            var port = Config.It.main.GetInt("main", "port", 30021);
            client.Start(ip, port);
        }

        public override void OnDestroy()
        {
            BindEvents(false);
        }

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;
            events.Bind(Card.EventDefine.ConnectFailed, onConnectFailed, bind);
            events.Bind(Card.EventDefine.ConnectionBroken, onConnectBroken, bind);
            events.Bind(Card.EventDefine.LoginSucc, onLoginSucc, bind);
        }

        private void onConnectFailed(params object[] args)
        {
            var e = args[0] as Card.HEventConnectFailed;

            UIMgr.It.GetPanel<PanelDialog>().ShowInfo($"连接失败: {e.error}");
        }

        private void onConnectBroken(params object[] args)
        {
            var e = args[0] as Card.HEventConnectBroken;

            UIMgr.It.GetPanel<PanelDialog>().ShowInfo($"连接断开: {e.error}");
        }

        private void onLoginSucc(params object[] args)
        {
            var e = args[0] as Card.HEventLoginSucc;                        
        }
    }
} // namespace Phoenix
