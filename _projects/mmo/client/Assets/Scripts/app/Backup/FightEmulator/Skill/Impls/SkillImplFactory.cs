using Phoenix.Core;
using System.Collections.Generic;


namespace Phoenix.Game.Skill
{	   
    public class SkillImplFactory : Singleton<SkillImplFactory>
    {
        private Dictionary<eSkillType, BaseSkillImpl> _impls =
            new Dictionary<eSkillType, BaseSkillImpl>();
        public SkillImplFactory()
        {
            registerAll();
        }

        private void registerAll()
        {
            register(eSkillType.Normal, new NormalSkillImpl());
            register(eSkillType.TargetBullet, new TargetBulletImpl());
        }

        private void register(eSkillType type, BaseSkillImpl impl)
        {
            _impls[type] = impl;
        }

        public void ApplyHit(eSkillType type, Skill skill)
        {
            BaseSkillImpl impl;
            if (!_impls.TryGetValue(type, out impl))
                return;
            impl.ApplyHit(skill);
        }
    }
} // namespace Phoenix
