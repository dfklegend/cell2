using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using UnityEngine.UI;
using Phoenix.Utils;

namespace Phoenix.Game.FightEmulator
{	
    public abstract class BaseModel
    {
        protected float _depth = 0;
        public float depth { get { return _depth; } }

        public abstract void SetPos(float x, float y);

        public abstract void Update();

        public abstract Transform GetTransform();
        public abstract void Destroy();
    }    
} // namespace Phoenix
