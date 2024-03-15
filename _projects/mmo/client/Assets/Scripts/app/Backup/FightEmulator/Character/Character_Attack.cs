using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    public partial class Character : IOwnerImpl, ICharacter
    {
        private float _nextAttackTime = 0;
        private int _enemy = -1;

        public void SetEnemy(int enemy)
        {
            _enemy = enemy;
        }

        private void autoAttack()
        {
            if (IsDead())
                return;


            if (_enemy < 0)
                return;

            
            tryAutoWeaponAttack(_enemy);
            // 自动释放技能
            _castCtrl.Update();

            var now = Time.time;
            if (now < _nextAttackTime)
                return;
            _nextAttackTime = now + 5.0f;
            //attackEnemy1(_enemy);
        }

        public void TestAttackEnemy(Skill.Skill skill, int tarId)
        {
            Character tar = FightCtrl.It.GetChar(tarId);
            if (tar == null)
                return;
            if (tar.IsDead())
                return;

            FormulaResult result = FormulaUtil.WeaponFormula(skill, this, tar) as FormulaResult;
            tar.ApplyDmg((int)result.data.Dmg);

            Log.LogCenter.Default.Debug($"{name}攻击{tar.name}，skill: {skill.skillId} 造成{(int)result.data.Dmg}点伤害");
            Log.LogCenter.Default.Debug($"{tar.name}剩余血量: {tar.charAttrs.GetHP()}");

            // 抛出事件
            Core.HEventUtil.Dispatch(Core.GlobalEvents.It.events,
                new HEventAttack(this, tar, result));
        }

        private void attackEnemy1(int tarId)
        {
            Character tar = FightCtrl.It.GetChar(tarId);
            if (tar == null)
                return;
            if (tar.IsDead())
                return;
            _skill.SkillToTar(tarId, "三连击", null);
        }

        private bool canAutoWeaponAttack()
        {
            // 主技能释放不能
            // 比如吟唱时候，是无法武器攻击的
            if (_skill.HasMainSkillRunning())
                return false;
            return true;
        }

        private void tryAutoWeaponAttack(int tarId)
        {
            Character tar = FightCtrl.It.GetChar(tarId);
            if (tar == null)
                return;
            if (tar.IsDead())
                return;
            if (!canAutoWeaponAttack())
                return;
            // 检查攻击范围
            if (CharUtil.DistTo(this, tar) > attackRange)
                return;

            var now = TimeUtil.Now();
            var weaponSkillId = Skill.SkillDefine.MeleeSkillId;
            if(this.attackRange > 2f)
                weaponSkillId = Skill.SkillDefine.RangeSkillId;
            if (_PACtrl.CanAttack(Skill.eHandType.MainHand, now))
            {
                var ext = new Skill.SkillArgsEx();
                ext.hand = Skill.eHandType.MainHand;
                _skill.BGSkillToTar(tarId, weaponSkillId, ext);
                _PACtrl.SureAttackDisplayDelay(Skill.eHandType.OffHand, now);
            }
            if(GetShieldBlock() == 0 && _PACtrl.CanAttack(Skill.eHandType.OffHand, now))
            {
                var ext = new Skill.SkillArgsEx();
                ext.hand = Skill.eHandType.OffHand;
                _skill.BGSkillToTar(tarId, weaponSkillId, ext);
                _PACtrl.SureAttackDisplayDelay(Skill.eHandType.MainHand, now);
            }
        }

        public void ApplyDmg(int dmg)
        {
            var hp = _charAttrs.GetHP();
            hp -= dmg;
            var dead = hp <= 0;
            if (hp < 0)
                hp = 0;
            _charAttrs.SetHP(hp);

            Core.HEventUtil.Dispatch(Core.GlobalEvents.It.events,
                new HEventHPChanged(this));

            if (dead)
                onDead();
        }

        private void onDead()
        {
            FightCtrl.It.GetWorld().DestroyEntity(_id);
        }

        public bool IsDead()
        {
            return _charAttrs.GetHP() == 0;
        }

        public void OnSkillStart(Skill.Skill skill)
        { 
            if(skill.GetHandType() == Skill.eHandType.MainHand ||
                skill.GetHandType() == Skill.eHandType.OffHand)
            {
                _PACtrl.OnAttackStart(skill.GetHandType(), TimeUtil.Now());
            }
        }

        public void OnSkillHit(Skill.Skill skill)
        {
            // push cd
            var cd = skill.cfgData.cd;
            if (cd <= 0)
                return;
            _cd.PushCD(skill.skillId, cd, TimeUtil.Now());
        }

        public void OnSkillStop(Skill.Skill skill, bool beBroken)
        {
            if (skill.GetHandType() == Skill.eHandType.MainHand ||
                skill.GetHandType() == Skill.eHandType.OffHand)
            {
                _PACtrl.OnAttackStop(skill.GetHandType());
            }
        }

        public void OnNormalAttack()
        {
        }

        public void CastSkillToTar(string skillId, int level)
        {
            int tarId = _enemy;
            var ext = new Skill.SkillArgsEx();
            ext.hand = Skill.eHandType.MainHand;

            _skill.SkillToTar(tarId, skillId, ext);
        }

        public bool HasMainSkillRunning()
        {
            return _skill.HasMainSkillRunning();
        }

        public bool IsCDOK(string skillId)
        {
            return _cd.IsCDOK(skillId, TimeUtil.Now());
        }

        public int GetShieldBlock()
        {
            var item = GetOffHandWeapon();
            if (item == null)
                return 0;
            var equip = item as EquipItem;
            return equip.GetBlock();            
        }
    }

}// namespace Phoenix
