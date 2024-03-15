using System.Collections.Generic;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    [StringType("PanelBagSelect")]
    public class PanelBagSelect : BasePanel
    {
        Button _btnClose;
        RectTransform _content;
        Text _title;

        // 显示viewData
        BagEnvData _data;
        List<BaseItemListStyle> _items = new List<BaseItemListStyle>();

        public override void OnReady()
        {
            SetDepth(PanelDepth.AboveNormal + 20);
            base.OnReady();

            _btnClose = TransformUtil.FindComponent<Button>(_root, "BG/btnClose");
            _btnClose.onClick.AddListener(onBtnClose);
            _title = TransformUtil.FindComponent<Text>(_root, "BG/title");

            _content = TransformUtil.FindComponent<RectTransform>(_root, "BG/scrollview/view/content");
        }

        private void onBtnClose()
        {
            Hide();
        }

        private void initItems()
        {
            destroyItems();
            buildItems();
        }

        private void destroyItems()
        {
            foreach (var item in _items)
                item.Destroy();
            _items.Clear();
        }

        private void buildItems()
        {
            ItemBagUtil.BuildItems(_items, _data, _content.transform);
        }

        protected override void onShow()
        {            
            initItems();
        }

        protected override void onHide()
        {
            destroyItems();
        }

        public void ShowBag(string title, BagEnvData data)
        {
            _data = data;
            _title.text = title;

            Show();
        }               
    }
} // namespace Phoenix
