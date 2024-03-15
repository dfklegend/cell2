using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    // 冒字
    [StringType("PanelFloatText")]
    public class PanelFloatText : BasePanel
    {
        Transform _textRoot;       
        public override void OnReady()
        {
            SetDepth(PanelDepth.AboveNormal);
            base.OnReady();

            _textRoot = _root.Find("root");
            FText.FloatTextMgr.It.Init(_textRoot.gameObject);

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

        public override void Update() 
        {
            FText.FloatTextMgr.It.Update();
        }
    }
} // namespace Phoenix
