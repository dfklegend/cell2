using Phoenix.Utils;
using System.Threading.Tasks;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    [StringType("PanelEnterScene")]
    public class PanelEnterScene : BasePanel
    {
        private Text _info;
        public override void OnReady()
        {
            SetDepth(PanelDepth.Highest);
            base.OnReady();

            _info = TransformUtil.FindComponent<Text>(_root, "BG/Content");
            BindEvents(true);            
        }

        public override void OnDestroy()
        {
            BindEvents(false);
        }        

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;
            events.Bind(Card.EventDefine.RefreshSceneInfo, onEnterScene, bind);
        }      
        
        private void setInfo(string info)
        {
            _info.text = info;
        }

        private async void onEnterScene(params object[] args)
        {
            Show();

            var cfgId = Card.DataCenter.It.SceneCfgId;
            
            //Card.S
            setInfo($"½øÈë³¡¾°: {Card.DataCenter.It.SceneServer}-{Card.DataCenter.It.SceneId}-{Card.DataCenter.It.SceneCfgId}");
            await Task.Delay(1500);
            Hide();
        }
    }
} // namespace Phoenix
