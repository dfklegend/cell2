using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using Phoenix.Fsm;
using Phoenix.Utils;

namespace Phoenix.Game.FightEmulator
{
    [IntType((int)eAIState.RandMove)]
    public class RandMove : BaseAIState
    {
        float _nextMove;
        public override void OnEnter()
        {
        }

        public override void OnLeave()
        {
        }

        public override void Update()
        {
            var now = Time.time;
            if (now < _nextMove)
                return;
            _nextMove = now + 1.0f;

            var owner = _ctrl.owner;
            var pos = owner.pos;
            _ctrl.owner.moveCtrl.MoveTo(pos.x + MathUtil.RandomF(-5f, 5f),
                pos.z + +MathUtil.RandomF(-5f, 5f));
        }

        public override bool IsOver() { return false; }
    }    
}// namespace Phoenix
