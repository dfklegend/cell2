using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    [StringType("PanelMain")]
    public class PanelMain : BasePanel
    {
        public override void OnReady()
        {
            SetDepth(PanelDepth.AboveNormal + 10);
            base.OnReady();            

            BindEvents(true);

            var btn = TransformUtil.FindComponent<Button>(_root, "BG/bottom/btnRole");
            btn.onClick.AddListener(onRole);

            btn = TransformUtil.FindComponent<Button>(_root, "BG/bottom/btnFight");
            btn.onClick.AddListener(onFight);

            btn = TransformUtil.FindComponent<Button>(_root, "BG/bottom/btnBag");
            btn.onClick.AddListener(onBag);
        }

        public override void OnDestroy()
        {
            BindEvents(false);
        }        

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;            
        }     

        private void onRole()
        {
            UIMgr.It.OpenPanel("PanelCharInfo").Show();
            UIMgr.It.OpenPanel("PanelItemBag").Hide();
        }

        private void onFight()
        {
            UIMgr.It.OpenPanel("PanelCharInfo").Hide();
            UIMgr.It.OpenPanel("PanelItemBag").Hide();
            
        }

        private void onBag()
        {
            UIMgr.It.OpenPanel("PanelCharInfo").Hide();
            BagCtrl.It.ShowItemBag();
        }
    }
} // namespace Phoenix
