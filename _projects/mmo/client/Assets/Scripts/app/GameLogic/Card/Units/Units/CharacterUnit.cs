using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;

namespace Phoenix.Game.Card
{	
    public class CharacterUnit : BaseLogicUnit
    {
        protected Character _char;
        public Character character { get { return _char; } }
        public void SetCharacter(Character c) { _char = c; }

        public override void Update() 
        {
            _char.Update();
        }

        public override void Destroy() 
        {
            UnitModelMgr.It.DestroyModel(_char.id);
        }
    }
    
} // namespace Phoenix
