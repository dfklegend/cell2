using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Log;

namespace Phoenix.Game
{
    [StringType("PanelCreateCard")]
    public class PanelCreateCard : BasePanel
    {
        private InputField _input;
        public override void OnReady()
        {
            SetDepth(PanelDepth.Normal + 100);
            base.OnReady();

            var btnCancel = _root.Find("BG/btnCancel").GetComponent<Button>();
            btnCancel.onClick.AddListener(() => {
                Hide();
            });

            var btnOK = _root.Find("BG/btnOK").GetComponent<Button>();
            btnOK.onClick.AddListener(() => {
                onOK();
            });

            _input = _root.Find("BG/inputName").GetComponent<InputField>();
        }                

        private void onOK()
        {
            var name = _input.text;
            if (name == "")
            {               
                UIMgr.It.GetPanel<PanelDialog>().ShowInfo("ÇëÊäÈëÃû×Ö");
                return;
            }
            GMCmdMgr.It.Execute($"lcmd createcard {name}", (succ) =>{
                if (!succ)
                    Debug.Log("create card failed!");
            });
            _input.text = "";
            Hide();
        }
    }
} // namespace Phoenix
