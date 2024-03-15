using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    [StringType("PanelFightChars")]
    public class PanelFightChars : BasePanel
    {
        Transform _charRoot;
        Transform _gridRoot;
        Text _logs;
        ScrollRect _scrollRect;       

        public override void OnReady()
        {
            SetDepth(100);
            base.OnReady();

            _charRoot = _root.Find("BG/chars");
            _gridRoot = _root.Find("BG/grid");
            _scrollRect = TransformUtil.FindComponent<ScrollRect>(_root, "BG/logs/scroll");
            _logs = TransformUtil.FindComponent<Text>(_root, "BG/logs/scroll/Text"); 

            _logs.text = "";

            BindEvents(true);
        }

        public override void OnDestroy()
        {
            BindEvents(false);
        }

        public Transform GetCharRoot()
        {
            return _charRoot;
        }

        public Transform GetGridRoot()
        {
            return _gridRoot;
        }

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;            
        }     
    }
} // namespace Phoenix
