using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using UnityEngine.UI;
using Phoenix.Utils;

namespace Phoenix.Game.Card
{	
    public class StaticModel : BaseModel
    {
        private GameObject _go;
        private Transform _root;
        private RectTransform _rect;
        private StaticUnit _owner;

        public void Init(GameObject go, StaticUnit unit)
        {
            if (!go)
                return;
            _go = go;
            _owner = unit;
            _root = go.transform;
            _rect = go.GetComponent<RectTransform>();          
        }

        public override void SetPos(float x, float y)
        {
            _rect.anchoredPosition = new Vector2(x, y);
            updateDepth(new Vector2(x, y));
        }

        public override void Update() 
        {
            UpdatePos();
        }

        public override Transform GetTransform() { return _root; }    
        

        public override void Destroy()
        {
            if (!_go)
                return;
            GameObject.Destroy(_go);
        }

        public void UpdatePos()
        {
            // 根据位置刷新界面位置
            Vector3 pos = _owner.pos;

            var uiPos = DisplayWorld.It.LogicToRelativeDisplay(pos);
            SetPos(uiPos.x, uiPos.y);            
        }
    }    
} // namespace Phoenix
