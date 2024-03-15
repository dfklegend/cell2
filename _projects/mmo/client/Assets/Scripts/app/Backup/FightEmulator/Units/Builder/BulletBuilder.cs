using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{	
    public class BulletCreateInfo
    {
        public string skillId;
        public float speed = 2f;
        public int tarId = -1;
        public int ownerId = -1;
    }

    public class BulletBuilder : Singleton<BulletBuilder>
    {       
        public Entity.Entity CreateBullet(EntityWorld world,
            BulletCreateInfo info)
        {
            var e = world.CreateEntity() as Entity.Entity;
            if (e == null)
                return null;

            var u = new BulletUnit();
            u.SetUnitType((int)eUnitType.Bullet);
            e.BindLogicUnit(u);

            var bullet = new Bullet(e.ID, info.ownerId,
                info.tarId, info.skillId, info.speed);
            u.SetBullet(bullet);

            UnitModelMgr.It.CreateBullet(u);

            return e;
        }       
    }
} // namespace Phoenix
