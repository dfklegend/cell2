using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using UnityEngine.UI;
using Phoenix.Utils;

namespace Phoenix.Game.Card
{	
    public class BulletModel : BaseModel
    {
        private GameObject _go;
        private Transform _root;
        private RectTransform _rect;
        //private Bullet _owner;

        public void Init(GameObject go, BulletUnit unit)
        {
            if (!go)
                return;
            _go = go;
            _root = go.transform;
            _rect = go.GetComponent<RectTransform>();

            initCtrls();
            setInfo(unit);
        }

        public override void SetPos(float x, float y)
        {
            _rect.anchoredPosition = new Vector2(x, y);
        }

        public override void Update() 
        {
            UpdatePos();
        }

        public override Transform GetTransform() { return _root; }

        private void initCtrls()
        {   
        }

        private void setInfo(BulletUnit unit)
        {
            //_owner = unit.bullet;
            UpdatePos();
        }

        public void UpdatePos()
        {
            // 根据位置刷新界面位置
            //Vector3 pos = _owner.pos;

            //var uiPos = DisplayWorld.It.WindowLogicToDisplay(pos);
            //SetPos(uiPos.x, uiPos.y);
            //updateDepth();
        }

        public override void Destroy()
        {
            if (!_go)
                return;
            GameObject.Destroy(_go);
        }
    }    
} // namespace Phoenix
