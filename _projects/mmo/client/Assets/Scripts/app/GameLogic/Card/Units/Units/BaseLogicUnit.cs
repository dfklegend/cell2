using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;

namespace Phoenix.Game.Card
{	
    public class BaseLogicUnit : ILogicUnit
    {
        protected Entity.Entity _entity;
        public Entity.Entity entity { get { return _entity; } }        

        public void SetEntity(Entity.IEntity e) 
        { 
            _entity = e as Entity.Entity; 
        }

        protected int _unitType = 0;
        public int unitType { get { return _unitType; } }

        public void SetUnitType(int ut) { _unitType = ut; }
        public int GetUnitType() { return _unitType; }        

        public virtual void Update() 
        {            
        }

        public virtual void Destroy() 
        {   
        }
    }
    
} // namespace Phoenix
