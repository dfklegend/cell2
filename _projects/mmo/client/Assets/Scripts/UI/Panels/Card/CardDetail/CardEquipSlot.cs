using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Log;

namespace Phoenix.Game
{
    // װ���ͼ��ܵ�װ��λ��
    public class CardEquipSlot
    {
        public delegate void OnClickSlot(CardEquipSlot box);

        // �����װ����λ
        int _slotIndex;
        IShowItem _item;

        Transform _root;       
        OnClickSlot _handler;

        BaseItemIconStyle _style;

        public void Init(Transform root,
            int index, OnClickSlot clickHandler)
        {
            _slotIndex = index;
            if (!root)
                return;
            _root = root;
            _handler = clickHandler;
            RefreshInfo(null);
        }

        public int GetIndex()
        {
            return _slotIndex;
        }

        public void RefreshInfo(IShowItem v)
        {
            _item = v;
            if (_style == null)
                _style = ItemIconStyleBuilder.CreateItemStyle(eItemStyle.Normal, _root, v, onClickBtn, null);
            else
                _style.RefreshInfo(v, null);
        }               

        private void onClickBtn(IShowItem v, object arg)
        {
            if (_handler != null)
                _handler(this);
        }
    }

} // namespace Phoenix
