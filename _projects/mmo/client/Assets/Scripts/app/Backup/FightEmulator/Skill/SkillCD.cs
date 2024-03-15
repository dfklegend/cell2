using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

using Phoenix.Game.FightEmulator;


namespace Phoenix.Game.Skill
{	    
    public class SkillCD
    {
        public string skillId;                
        public float endTime;
        public float getRestTime(float now)
        {
            if (endTime <= 0)
                return 0;
            var result = endTime - now;
            if (result < 0)
                result = 0;
            return result;
        }

        public bool IsCDOK(float now)
        {
            return getRestTime(now) <= 0;
        }
    }

    public class SkillCDCtrl
    {
        private Dictionary<string, SkillCD> _cds
            = new Dictionary<string, SkillCD>();

        private SkillCD getObj(string skillId)
        {
            SkillCD obj;
            _cds.TryGetValue(skillId, out obj);
            return obj;
        }

        private void addObj(SkillCD obj)
        {
            _cds[obj.skillId] = obj;
        }

        public bool IsCDOK(string skillId, float now)
        {
            var obj = getObj(skillId);
            if (obj == null)
                return true;
            return obj.IsCDOK(now);
        }

        public void PushCD(string skillId, float cd, float now)
        {
            if (cd <= 0)
                return;
            var obj = getObj(skillId);
            if(obj == null)
            {
                obj = new SkillCD();
                obj.skillId = skillId;
                obj.endTime = now + cd;
                addObj(obj);
                return;
            }

            if(obj.endTime < now + cd)
            {
                obj.endTime = now + cd;
            }
        }

        public void SubCD(string skillId, float offset)
        {
            if (offset < 0)
                return;
            var obj = getObj(skillId);
            if (obj == null)
                return;
            obj.endTime -= offset;
        }

        public void ClearCD(string skillId)
        {
            var obj = getObj(skillId);
            if (obj == null)
                return;
            obj.endTime = 0;
        }
    }
} // namespace Phoenix
