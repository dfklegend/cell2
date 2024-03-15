using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Log;

namespace Phoenix.Game
{
    [StringType("PanelDialog")]
    public class PanelDialog : BasePanel
    {   
        private Button _btnOK;
        private Text _content;
        public override void OnReady()
        {
            SetDepth(PanelDepth.Highest);
            base.OnReady();

            _btnOK = _root.Find("BG/btnOK").GetComponent<Button>();
            _btnOK.onClick.AddListener(() => {
                Hide();
            });

            _content = _root.Find("BG/Content").GetComponent<Text>();
        }        

        public void ShowInfo(string str)
        {
            _content.text = str;
            Show();
        }
    }
} // namespace Phoenix
