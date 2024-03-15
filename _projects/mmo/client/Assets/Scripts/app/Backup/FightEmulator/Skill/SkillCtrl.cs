using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

using Phoenix.Game.FightEmulator;


namespace Phoenix.Game.Skill
{	    
    // 一些额外参数
    public class SkillArgsEx
    {
        // 那个手触发的技能，主要是自动攻击
        public eHandType hand = eHandType.None;
        // 是否只是结算
        public bool onlyBalance = false;
    }

    public class SkillCtrl
    {
        private Character _owner;
        // 主技能相互互斥
        private Skill _mainSkill;
        // 后台技能
        private List<Skill> _bgSkills = new List<Skill>();

        public void Init(Character owner)
        {
            _owner = owner;
        }

        public void Update()
        {
            updateMainSkill();
            updateBGSkills();
        }

        private void updateMainSkill()
        {
            if (_mainSkill == null)
                return;
            if(_mainSkill.IsOver())
            {
                _mainSkill = null;
            }
            else
            {
                _mainSkill.Update();
            }
        }

        private void updateBGSkills()
        {
            for(var i = 0; i < _bgSkills.Count;)
            {
                var skill = _bgSkills[i];
                if(skill.IsOver()) 
                {
                    _bgSkills.RemoveAt(i);
                }
                else
                {
                    skill.Update();
                    i++;
                }
            }
        }        

        public Skill SkillToTar(int tarId, string skillId, SkillArgsEx ext)
        {
            if(HasMainSkillRunning())
            {
                // TODO: 技能顶替
                return null;
            }
            var skill = startSkillToTar(tarId, skillId, ext);
            if(skill != null)                
                _mainSkill = skill;
            return skill;
        }

        public bool HasMainSkillRunning()
        {
            return _mainSkill != null;
        }

        public Skill BGSkillToTar(int tarId, string skillId, SkillArgsEx ext)
        {   
            var skill = startSkillToTar(tarId, skillId, ext);
            if (skill != null)
                _bgSkills.Add(skill);            
            return skill;
        }

        private Skill startSkillToTar(int tarId, string skillId, SkillArgsEx ext)
        {
            var skill = new Skill();
            if (!skill.Init(_owner, skillId))
                return null;
            
            skill.SetTarId(tarId);
            // set ext
            if (ext != null)
            {
                skill.SetHandType(ext.hand);
                skill.SetOnlyBalance(ext.onlyBalance);
            }

            skill.Start();
            return skill;
        }
    }
} // namespace Phoenix
