using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{ 
    // 用于做目标挑选
    public class SearchCtrl : Singleton<SearchCtrl>
    {
        EntityWorld _world;
        public void Init(EntityWorld world)
        {
            _world = world;
        }

        public void Clear()
        {
            _world = null;
        }

        public Character FindEnemy(Character src)
        {
            var visitor = new IsEnemy(src);
            _world.Visit(visitor);
            return visitor.found;
        }

        public Character FindNearestEnemy(Character src,
            float range)
        {
            var visitor = new FindNearestEnemy(src, range);
            _world.Visit(visitor);
            return visitor.found;            
        }
    }

} // namespace Phoenix
