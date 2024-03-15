using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using Phoenix.Fsm;
using Phoenix.Utils;

namespace Phoenix.Game.FightEmulator
{
    [IntType((int)eAIState.Attack)]
    public class Attack : BaseAIState
    {   
        public override void OnEnter()
        {
            _ctrl.owner.moveCtrl.Stop();
        }

        public override void OnLeave()
        {
        }

        // TODO:
        // 根据当前选择技能来移动
        // 技能cast会需要等待一阵(等待移动满足施法距离)
        public override void Update()
        {
            var ai = _ctrl.ai;
            var owner = _ctrl.owner;
            Character tar = FightCtrl.It.GetChar(ai.enemyId);
            if (tar == null)
                return;
            float attackRange = owner.attackRange;
            float dist = CharUtil.DistTo(owner, tar);
            
            if(dist > attackRange)
            {
                var tarPos = CharUtil.GetMoveTar(owner, tar, attackRange*0.8f);
                owner.moveCtrl.MoveTo(tarPos.x, tarPos.z);
            }
            // 太近
            if(dist <= attackRange)
            {
                owner.moveCtrl.Stop();
            }
        }

        public override bool IsOver() 
        {
            var ai = _ctrl.ai;
            if (!IsTarValid(ai.enemyId))
                return true;
            return false; 
        }

        private bool IsTarValid(int id)
        {
            Character tar = FightCtrl.It.GetChar(id);
            if (tar == null || tar.IsDead())
                return false;
            return true;
        }
    }    
}// namespace Phoenix
