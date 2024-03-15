using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;

namespace Phoenix.Game.Card
{	
    // ¾²Ì¬¶ÔÏó
    public class StaticUnit : BaseLogicUnit
    {
        protected Vector3 _pos;
        public Vector3 pos => _pos;
        public void SetPos(Vector3 pos)
        {
            _pos = pos;
        }
    }
    
} // namespace Phoenix
