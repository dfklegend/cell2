using UnityEngine;

namespace Phoenix.Game
{
    // . 面板异步载入
    // . 使用事件来与面板交互
    public abstract class BasePanel
    {
        private bool _loading = false;
        private GameObject _go;
        protected Transform _root;
        private bool _destroyed = false;
        public bool Destroyed { get { return _destroyed; } }

        public bool autoRecycle = false;
        private bool _visible = true;

        // 面板顺序
        // 用于排序 越大越在前面
        private int _depth = PanelDepth.Normal;
        public int depth { get { return _depth; } }

        public GameObject GetGameObject() { return _go; }
        public bool IsLoading() { return _loading; }
        public void SetLoading() { _loading = true; }

        public bool IsReady()
        {
            return _go != null;
        }

        // 面板创建
        public virtual void OnCreate() { }
        public void OnLoaded(GameObject go)
        {
            if (!go)
                return;
            if (_destroyed)
            {
                GameObject.Destroy(go);
                return;
            }
            _go = go;
            _root = go.transform;

            UIMgr.It.AddToScene(this);            
            OnReady();
            applyVisible();
        }

        // 面板ready
        // 可以获取控件做操作
        public virtual void OnReady()
        {

        }

        public void Destroy()
        {
            _destroyed = true;
            if (_go)
            {
                UIMgr.It.RemoveFromScene(this);
                GameObject.Destroy(_go);
                OnDestroy();
                _root = null;
            }
        }

        public virtual void OnDestroy() { }


        public void Show()
        {
            _visible = true;
            applyVisible();
        }

        public void Hide()
        {
            _visible = false;
            applyVisible();
        }

        private void applyVisible()
        {
            if (!_go)
                return;
            _go.SetActive(_visible);
            
            if (_visible)
                onShow();
            else
                onHide();
        }

        public bool IsVisible()
        {
            return _visible && IsReady();
        }
        

        protected virtual void onShow() {}
        protected virtual void onHide() {}

        public void SetDepth(int depth)
        {
            _depth = depth;
            UIMgr.It.SetZDirt();
        }

        public virtual void Update() { }
    }
} // namespace Phoenix
