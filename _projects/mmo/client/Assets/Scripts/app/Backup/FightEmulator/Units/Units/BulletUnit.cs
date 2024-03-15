using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;

namespace Phoenix.Game.FightEmulator
{	
    public class BulletUnit : BaseLogicUnit
    {
        protected Bullet _bullet;
        public Bullet bullet { get { return _bullet; } }
        public void SetBullet(Bullet v) { _bullet = v; }

        public override void Update() 
        {
            _bullet.Update();
            if(_bullet.IsOver())
            {
                FightCtrl.It.GetWorld().DestroyEntity(bullet.id);
            }
        }

        public override void Destroy() 
        {
            UnitModelMgr.It.DestroyBaseModel(_bullet.id);
        }
    }
    
} // namespace Phoenix
