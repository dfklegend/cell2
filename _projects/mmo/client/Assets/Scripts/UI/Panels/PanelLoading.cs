using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    [StringType("PanelLoading")]
    public class PanelLoading : BasePanel
    {
        public override void OnReady()
        {
            SetDepth(PanelDepth.Highest);
            base.OnReady();            

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
    }
} // namespace Phoenix
