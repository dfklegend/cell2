using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

using Phoenix.Game.FightEmulator;


namespace Phoenix.Game.Skill
{	   
    // 对目标发射一个子弹
    public class TargetBulletImpl : BaseSkillImpl
    {
        public override void ApplyHit(Skill skill)
        {
            BulletCreateInfo info = new BulletCreateInfo();
            info.skillId = skill.skillId;
            info.ownerId = skill.owner.id;
            info.tarId = skill.tarId;
            info.speed = 5f;

            BulletBuilder.It.CreateBullet(FightCtrl.It.GetWorld(), info);
        }
    }
} // namespace Phoenix
