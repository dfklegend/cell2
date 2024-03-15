using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator.BagSystem;

namespace Phoenix.Game
{   
    // 抽象物品显示
    public class NormalIconStyle : BaseItemIconStyle
    {
        Image _icon;
        Text _name;

        public override void OnReady()
        {
            _icon = TransformUtil.FindComponent<Image>(_root, "icon");
            _name = TransformUtil.FindComponent<Text>(_root, "name");
            TransformUtil.FindComponent<Button>(_root, "btn").onClick.AddListener(onClick);
        }

        protected override void onRefresh(IShowItem v, IIconStyleOptions options)
        {
            if (v == null)
            {
                reset();
                return;
            }            
            
            if(options != null &&  options.NeedHideName())
            {
                _name.gameObject.SetActive(false);
            }
            else
                _name.text = v.GetName();            
            _icon.sprite = UIUtil.LoadIcon(v.GetIcon());
        }

        private void reset()
        {
            _icon.sprite = null;
            _name.text = "";
        }

        protected void onClick()
        {
            invoke();
        }
    }
} // namespace Phoenix
