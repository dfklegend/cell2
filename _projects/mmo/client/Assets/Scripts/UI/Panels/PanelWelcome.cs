using Phoenix.Utils;
using UnityEngine;

namespace Phoenix.Game
{
    [StringType("PanelWelcome")]
    public class PanelWelcome : BasePanel
    {
        float _timeReady = 0.0f;
        public override void OnReady()
        {
            SetDepth(200);
            base.OnReady();

            _timeReady = Time.time;
        }

        public override void Update()
        {            
            if(_timeReady > 0 && Time.time > _timeReady + 2.0)
            {
                ClientApp.It.stateCtrl.ChangeState((int)eAppState.Login);
                Destroy();
            }
        }
    }
} // namespace Phoenix
