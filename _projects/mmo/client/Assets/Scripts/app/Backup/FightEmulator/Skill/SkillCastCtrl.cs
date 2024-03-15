using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

using Phoenix.Game.FightEmulator;


namespace Phoenix.Game.Skill
{	    
    // 依次释放技能
    public class SkillCastCtrl
    {
        private SkillSlots _slots;
        private Character _owner;
        // 下次选择技能的时间
        private float _nextSelectSkill;
        private int _indexSlot = 0;

        public void Init(Character c, SkillSlots slots)
        {
            _owner = c;
            _slots = slots;
            _nextSelectSkill = TimeUtil.Now() + 3f;         
        }

        public void Update()
        {
            var now = TimeUtil.Now();
            if(now >= _nextSelectSkill)
            {
                selectAndCast();
                _nextSelectSkill = now + 1f;
            }
        }

        private void selectAndCast()
        {
            if (hasMainSkillRunning())
                return;
            if (tryCastSuperSkill())
                return;

            int tryTimes = 0;
            while(tryTimes < _slots.GetValidNormalNum())
            {
                var slot = _slots.GetSlot(_indexSlot);
                if (!slot.IsValid() || !canCast(slot))
                {
                    goNextSlot();
                    tryTimes++;
                    continue;
                }

                // cast it
                castSkill(slot);
                goNextSlot();
                break;
            }
        }

        private void goNextSlot()
        {
            _indexSlot++;
            if (_indexSlot >= SkillDefine.IndexSuperSkill)
                _indexSlot = 0;
        }

        private bool isPassive(SkillEntry slot)
        {
            return false;
        }

        private bool canCast(SkillEntry slot)
        {
            if (isPassive(slot) ||
                !_owner.IsCDOK(slot.skillId))
                return false;
            return true;
        }

        //private bool

        private bool hasMainSkillRunning()
        {
            return _owner.HasMainSkillRunning();
        }

        private bool tryCastSuperSkill()
        {
            // 是否能量满了
            return false;
        }

        // TODO: 技能释放距离添加后
        // 需要添加延迟施法尝试(距离可能不够)
        private void castSkill(SkillEntry slot)
        {
            // TODO:根据技能特性选择目标释放
            // 有治疗技能再调试
            _owner.CastSkillToTar(slot.skillId, slot.level);
            _owner.moveCtrl.Stop();
        }
    }
} // namespace Phoenix
