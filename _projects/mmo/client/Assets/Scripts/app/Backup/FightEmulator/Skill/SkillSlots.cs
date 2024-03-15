using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

using Phoenix.Game.FightEmulator;


namespace Phoenix.Game.Skill
{	    
    public class SkillEntry
    {
        public string skillId;
        public int level;
        public SkillData sd;

        public void Reset()
        {
            skillId = "";
            level = 0;
            sd = null;
        }

        public bool IsValid()
        {
            return level > 0;
        }
    }

    public class SkillSlots
    {
        private List<SkillEntry> _slots
            = new List<SkillEntry>();
        // 除超级技能外
        private int _validNormalNum;
        private bool _dirtNormalNum = false;

        public SkillSlots()
        {
            init();
        }

        private void init()
        {
            for(var i = 0; i < SkillDefine.MaxSkillSlot; i ++)
            {
                _slots.Add(new SkillEntry());
            }
        }

        public SkillEntry GetSlot(int index)
        {
            if (index < 0 || index >= SkillDefine.MaxSkillSlot)
                return null;
            return _slots[index];
        }

        public void SetSkill(int index, string skillId, int level)
        {
            var slot = GetSlot(index);
            if (slot == null)
                return;
            if (-1 != GetSkillIndex(skillId))
                return;
            var sd = SkillDataMgr.It.GetItem(skillId);
            if (sd == null)
                return;

            if(index == SkillDefine.IndexSuperSkill)
            {
                // 只能装备超级技能
            }

            slot.skillId = skillId;
            slot.level = level;
            slot.sd = sd;

            _dirtNormalNum = true;
        }

        public int GetSkillIndex(string skillId)
        {
            for(var i = 0; i < SkillDefine.MaxSkillSlot; i ++)
            {
                if (skillId == _slots[i].skillId)
                    return i;
            }
            return -1;
        }

        public void RemoveSkill(int index)
        {
            var slot = GetSlot(index);
            if (slot == null)
                return;
            slot.Reset();
        }

        public void updateValidNormalNum()
        {
            _validNormalNum = 0;
            for(var i = 0; i < SkillDefine.IndexSuperSkill; i ++)
            {
                if (_slots[i].IsValid())
                    _validNormalNum++;
            }
        }

        public int GetValidNormalNum()
        { 
            if(_dirtNormalNum)
            {
                updateValidNormalNum();
                _dirtNormalNum = false;
            }
            return _validNormalNum;
        }
    }
} // namespace Phoenix
