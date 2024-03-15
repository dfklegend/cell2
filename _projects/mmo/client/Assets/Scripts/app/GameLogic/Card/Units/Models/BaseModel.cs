using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using UnityEngine.UI;
using Phoenix.Utils;

namespace Phoenix.Game.Card
{	
    public abstract class BaseModel
    {
        // 显示排序深度，越大越在外面
        protected float _depth = 0;
        public float depth { get { return _depth; } }

        public abstract void SetPos(float x, float y);

        public abstract void Update();

        public abstract Transform GetTransform();
        public abstract void Destroy();

        protected void updateDepth(Vector2 pos)
        {            
            // Y轴负向            
            var factor = -pos.y * 1f + pos.x / 1000;
            _depth = factor;
        }
    }    
} // namespace Phoenix
