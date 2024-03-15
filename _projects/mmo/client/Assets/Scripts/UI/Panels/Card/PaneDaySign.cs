using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Core;
using Phoenix.Log;
using Network;
using System.Text;

namespace Phoenix.Game
{
    [StringType("PanelDaySign")]
    public class PanelDaySign : BasePanel
    {   
        private Button _btnClose;       
        private Button _btnSign;       
        private Text _textSigned;
      
        public override void OnReady()
        {
            SetDepth(1000);
            base.OnReady();

            _btnClose = _root.Find("BG/btnClose").GetComponent<Button>();
            _btnClose.onClick.AddListener(() => {
                Hide();
            });


            _btnSign = TransformUtil.FindComponent<Button>(_root, "BG/bottom/btnSign");
            _btnSign.onClick.AddListener(() => {
               onBtnSign();
            });
            _textSigned = TransformUtil.FindComponent<Text>(_root, "BG/bottom/signed");
            

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

        protected override void onShow()
        {
            base.onShow();
            var system = Singleton<Systems>.It.GetSystem<DaySignSystem>("daysign");
            updateState(system.signed);
        }

        private void updateState(bool signed)                 
        {
            if (signed)
            {
                _btnSign.gameObject.SetActive(false);
                _textSigned.gameObject.SetActive(true);              
            }
            else
            {
                _btnSign.gameObject.SetActive(true);
                _textSigned.gameObject.SetActive(false);
            }            
        }

        private void onBtnSign() 
        {
            var system = Singleton<Systems>.It.GetSystem<DaySignSystem>("daysign");
            system.Sign((signed)=> {
                updateState(signed);
            });
        }
    }
} // namespace Phoenix
