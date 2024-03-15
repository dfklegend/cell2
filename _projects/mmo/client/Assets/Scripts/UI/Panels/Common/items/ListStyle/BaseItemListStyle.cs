using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator.BagSystem;

namespace Phoenix.Game
{    
    // 抽象物品显示
    public abstract class BaseItemListStyle : IShowItemStyle
    {
        protected GameObject _go;
        protected Transform _root;
        public IShowItemHandler handler;

        public void Init(GameObject go)
        {
            _go = go;
            _root = go.transform;
            OnReady();
        }
                
        public abstract void OnReady();
        public abstract void RefreshInfo(IShowItem item);

        // call by ItemStyleHolder
        public virtual void OnClick(IShowItem item) 
        {
            handler?.OnClick(this, item);
        }

        public void Destroy()
        {
            if (!_go)
                return;
            GameObject.Destroy(_go);
        }
    }
} // namespace Phoenix
