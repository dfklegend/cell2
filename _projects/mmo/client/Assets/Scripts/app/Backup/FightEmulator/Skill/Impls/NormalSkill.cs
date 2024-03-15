using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;



namespace Phoenix.Game.Skill
{	   
    public class NormalSkillImpl : BaseSkillImpl
    {
        public override void ApplyHit(Skill skill)
        {
            // TODO
            // 首先根据技能类型收集有效目标
            // 然后对每一个目标执行结算
            skill.ApplyOneTarget(skill.tarId);
        }
    }
} // namespace Phoenix
