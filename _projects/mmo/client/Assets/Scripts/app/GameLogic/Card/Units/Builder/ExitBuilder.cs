using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.Card
{
    public class ExitCreateInfo
    {
        public Vector3 pos;
    }

    public class ExitBuilder : Singleton<ExitBuilder>
    {
        public Entity.Entity CreateExit(EntityWorld world,
            ExitCreateInfo info)
        {
            var e = world.CreateEntity() as Entity.Entity;
            if (e == null)
                return null;
            Log.LogCenter.Default.Debug("CreateExit");
            var u = new StaticUnit();
            u.SetUnitType((int)eUnitType.Exit);
            e.BindLogicUnit(u);

            u.SetPos(info.pos);
            UnitModelMgr.It.CreateExit(u);
            
            return e;
        }
    }
} // namespace Phoenix
