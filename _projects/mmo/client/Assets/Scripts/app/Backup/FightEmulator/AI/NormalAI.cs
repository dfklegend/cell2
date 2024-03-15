using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using Phoenix.Fsm;
using Phoenix.Utils;

namespace Phoenix.Game.FightEmulator
{
    public partial class AIFunc
    {
        // 有敌人就切成攻击
        public static void NormalUpdate(IFSMCtrl inCtrl)
        {
            var ctrl = inCtrl as AIStateCtrl;
            var ai = ctrl.ai;
            var owner = ctrl.owner;

            if(ctrl.GetStateType() == (int)eAIState.Init)
            {
                ctrl.ChangeState((int)eAIState.RandMove);
                return;
            }

            if(ctrl.GetStateType() == (int)eAIState.RandMove)
            {
                // 尝试索敌
                var found = SearchCtrl.It.FindNearestEnemy(owner, owner.guardRange);
                if (found != null)
                {
                    Log.LogCenter.Default.Debug($"{owner.id} found enemy {found.id}");
                    ai.SetEnemy(found);
                    ctrl.ChangeState((int)eAIState.Attack);
                }
            }

            if (ctrl.GetStateType() == (int)eAIState.Attack)
            {
                if (ctrl.GetState().IsOver())
                {
                    ctrl.ChangeState((int)eAIState.RandMove);
                    Log.LogCenter.Default.Debug("AttackState Over");
                }
            }
        }
    }
}// namespace Phoenix
