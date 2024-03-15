using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator.BagSystem;
using System;

namespace Phoenix.Game
{   
    // 抽象物品显示
    public abstract class BaseItemIconStyle
    {
        protected GameObject _go;
        protected Transform _root;
        
        // 点击之后的响应
        protected Action<IShowItem, object> _handler;
        protected object _arg;
        protected IShowItem _item;

        public void Init(GameObject go)
        {
            _go = go;
            _root = go.transform;
            OnReady();
        }

        public void Destroy()
        {
            if (!_go)
                return;
            GameObject.Destroy(_go);
            _go = null;
        }
                
        public abstract void OnReady();
        public void RefreshInfo(IShowItem item, IIconStyleOptions options)
        {
            _item = item;
            onRefresh(item, options);
        }

        protected abstract void onRefresh(IShowItem item, IIconStyleOptions options);
        
        public void SetHandler(Action<IShowItem, object> handler, object arg)
        {
            _handler = handler;
            _arg = arg;
        }

        protected void invoke()
        {
            if (_handler == null)
                return;
            _handler(_item, _arg);
        }        
    }
} // namespace Phoenix
