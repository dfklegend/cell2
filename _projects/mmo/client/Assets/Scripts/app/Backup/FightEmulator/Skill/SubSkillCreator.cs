using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

using Phoenix.Game.FightEmulator;


namespace Phoenix.Game.Skill
{	    
    public static class SubSkillCreator
    {
        public static Skill CreateSubSkill(Character owner,
            Skill parent, string subSkillId, int depth)
        {
            Skill sub = new Skill();
            if (!sub.Init(owner, subSkillId))
                return null;
            // copy目标信息
            sub.SetTarId(parent.tarId);
            sub.SetStackDepth(depth);
            return sub;
        }
    }
} // namespace Phoenix
