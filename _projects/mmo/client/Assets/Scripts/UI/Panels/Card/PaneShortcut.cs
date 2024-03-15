using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Log;
using Network;
using System.Text;

namespace Phoenix.Game
{
    [StringType("PanelShortcut")]
    public class PanelShortcut : BasePanel
    {   
        private Button _btnClose;       
      
        public override void OnReady()
        {
            SetDepth(1000);
            base.OnReady();

            _btnClose = _root.Find("BG/btnClose").GetComponent<Button>();
            _btnClose.onClick.AddListener(() => {
                Hide();
            });


            Button btn;
            string cmd;
            for ( var i = 0; i < 3; i ++)
            {
                btn = _root.Find("BG/btns/btn" + i).GetComponent<Button>();
                cmd = "lcmd enterscene " + (i + 1);
                btn.onClick.AddListener(createCmd(cmd));
            }

            btn = _root.Find("BG/btns/btn" + 3).GetComponent<Button>();
            btn.onClick.AddListener(() => {
                test();
                Hide();
            });

            BindEvents(true);
        }

        public override void OnDestroy()
        {
            BindEvents(false);
        }

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;  
        }
        
        private UnityEngine.Events.UnityAction createCmd(string cmd)
        {
            return () => { doCmd(cmd); };
        }

        private void doCmd(string cmd)
        {
            GMCmdMgr.It.Execute(cmd, null);
            Hide();
        }

        private void test()
        {
            doCmd("lcmd randswitchline");
            //var args = new Cproto.TestAdd()
            //{
            //    I = 10,
            //    J = 21
            //};

            //Systems.It.GetSystem("example").Request<Cproto.TestAddRet>("add", args, (ret, code) =>
            //{
            //    if (ret == null)
            //        return;
            //    GMCmdUtils.Output($"add result: {ret.Result} {code}");
            //});

            //var client = ClientApp.It.client;
            //client.SystemCmd<Cproto.TestAddRet>("example", "add", args, (ret, code) =>
            //{
            //    if (ret == null)
            //        return;
            //    GMCmdUtils.Output($"add result: {ret.Result} {code}");
            //});
        }
    }
} // namespace Phoenix
