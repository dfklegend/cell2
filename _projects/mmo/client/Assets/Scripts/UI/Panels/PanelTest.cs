using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    [StringType("PanelTest")]
    public class PanelTest : BasePanel
    {
        private Button _btnFightEmulator;
        private Button _btnOne;
        public override void OnReady()
        {
            SetDepth(100);
            base.OnReady();

            _btnFightEmulator = _root.Find("BG/btnFightEmulator").GetComponent<Button>();
            _btnFightEmulator.onClick.AddListener(() => {
                ClientApp.It.stateCtrl.ChangeState((int)eAppState.FightFE);
                Destroy();
            });

            _btnOne = _root.Find("BG/btnOne").GetComponent<Button>();
            _btnOne.onClick.AddListener(() => {
                //App.App.AppState.ChangeState((int)eAppState.Fight);
                Destroy();
            });
        }
    }
} // namespace Phoenix
