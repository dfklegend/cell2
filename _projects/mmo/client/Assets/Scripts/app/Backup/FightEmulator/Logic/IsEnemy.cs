using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{ 
    public class IsEnemy : IVisitor
    {
        Character _src;
        Character _found;

        public IsEnemy(Character src)
        {
            _src = src;
        }

        public Character found { get { return _found; } }

        public void Visit(IEntity e)
        {
            var tar = getCharacter(e);
            if (_src.side == tar.side)
                return;
            if (tar.IsDead())
                return;
            _found = tar;
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
