using System;
using System.Collections.Generic;
using Phoenix.Core;


namespace Phoenix.Game.FightEmulator
{
    public class ClassFeatureFactory : Singleton<ClassFeatureFactory>
    {
        private Dictionary<int, IAttrFeature> _features
            = new Dictionary<int, IAttrFeature>();
        private IAttrFeature _default;
        ClassFeatureFactory()
        {
            initFeatures();
        }

        private void initFeatures()
        {
            _features[(int)eClass.None] = new NoneFeature();
            _features[(int)eClass.Warrior] = new WarriorFeature();
            _features[(int)eClass.Rogue] = new RogueFeature();

            _default = _features[(int)eClass.None];
        }

        public IAttrFeature GetClassFeature(int classType)
        {
            IAttrFeature ret;
            if (_features.TryGetValue(classType, out ret))
                return ret;
            return _default;
        }
    }
}// namespace Phoenix
