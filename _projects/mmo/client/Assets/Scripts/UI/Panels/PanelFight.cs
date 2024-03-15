using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    [StringType("PanelFight")]
    public class PanelFight : BasePanel
    {
        private Button _btnExit;
        public override void OnReady()
        {
            SetDepth(100);
            base.OnReady();

            _btnExit = _root.Find("BG/btnExit").GetComponent<Button>();
            _btnExit.onClick.AddListener(() => {
                GameLogicCtrl.It.Clear();
                
                ClientApp.It.stateCtrl.ChangeState((int)eAppState.Login);
                Destroy();
            });
        }
    }
} // namespace Phoenix
