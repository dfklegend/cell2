using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{ 
    public class FindNearestEnemy : IVisitor
    {
        Character _src;
        float _range;
        
        Character _found;
        float _curDist = 0f;

        public FindNearestEnemy(Character src, float range)
        {
            _src = src;
            _range = range;
        }

        public Character found { get { return _found; } }

        public void Visit(IEntity e)
        {
            var tar = getCharacter(e);
            if (tar == null)
                return;
            if (_src.side == tar.side)
                return;
            if (tar.IsDead())
                return;
            float dist = CharUtil.DistTo(_src, tar);
            if (dist > _range)
                return;            
            if(_found == null || dist < _curDist)
            {
                _curDist = dist;
                _found = tar;                
            }            
        }

        Character getCharacter(IEntity eIn)
        {
            Entity.Entity e = eIn as Entity.Entity;
            CharacterUnit unit = e.unit as CharacterUnit;
            if (unit == null)
                return null;
            return unit.character;
        }
    }

} // namespace Phoenix
