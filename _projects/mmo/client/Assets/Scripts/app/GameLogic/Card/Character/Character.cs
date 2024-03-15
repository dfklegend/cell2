using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.Card
{
    // 可战斗主体
    public partial class Character
    {
        private string _name;
        public string name { get { return _name; } }
        public void SetName(string s)
        {
            _name = s;
        }

        private int _id = -1;
        public int id { get { return _id; } }
        public void SetId(int id) { _id = id; }

        // 场景id
        private int _sceneId = -1;
        public int sceneId { get { return _sceneId; } }
        public void SetSceneId(int id) { _sceneId = id; }

        // 位置
        public Vector3 pos = Vector3.zero;


        private int _level = 1;
        public int level { get { return _level; } }
        public void SetLevel(int level) { _level = level; }
        

        // 阵营
        private int _side = 0;
        public void SetSide(int side) { _side = side; }
        public int side { get { return _side; } }

        public float guardRange = 6.0f;
        public float attackRange = 2.0f;

        private int _hp = 100;
        private int _hpMax = 100;

        private Dictionary<int, AttrItem> _attrs = new Dictionary<int, AttrItem>();
        
        public MoveCtrl moveCtrl = new MoveCtrl();

        public Character()
        {
            moveCtrl.Init(this);
        }

        public int GetLevel()
        {
            return _level;
        }
        

        public void ChangeLevel(int level)
        {
            if (level == _level)
                return;           

            var off = level - _level;
            _level = level;            
        }        

        public float GetHPPercent()
        {
            return (float)_hp / _hpMax;
        }

        public float GetMPPercent()
        {
            return GetAttr(eAttrType.Energy) / 100;
        }

        public void SetHP(int hp)
        {
            this._hp = hp;
        }

        public int GetHP() 
        {
            return _hp;
        }

        public void SetHPMax(int max)
        {
            if (max <= 0)
                return;
            _hpMax = max;
        }

        public int GetHPMax()
        {
            return _hpMax;
        }

        public void SetAttr(int which, float value)
        {
            AttrItem item;
            if(_attrs.TryGetValue(which, out item))
            {
                item.value = value;
                return;
            }

            item = new AttrItem();
            item.index = which;
            item.value = value;
            _attrs[which] = item;
        }

        public float GetAttr(eAttrType et, float def = 0)
        {
            return GetAttr((int)et, def);
        }

        public float GetAttr(int which, float def = 0)            
        {
            AttrItem item;
            if (_attrs.TryGetValue(which, out item))
            {   
                return item.value;
            }
            return def;
        }

        public int GetIntAttr(int which, int def = 0)
        {
            return (int)GetAttr(which, def);
        }


        public void Update()
        {
            moveCtrl.Update();
        }       
    }
}// namespace Phoenix
