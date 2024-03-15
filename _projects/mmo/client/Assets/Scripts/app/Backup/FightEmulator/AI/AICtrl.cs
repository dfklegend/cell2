using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using Phoenix.Fsm;
using Phoenix.Utils;

namespace Phoenix.Game.FightEmulator
{
    public class AICtrl
    {
        private Character _owner;
        private AIStateCtrl _stateCtrl = AIStateCtrl.New();
        private int _enemyId = -1;
        public int enemyId { get { return _enemyId; } }
        

        public void Init(Character owner)
        {
            _owner = owner;
            _stateCtrl.SetAI(this);
            _stateCtrl.SetCharacter(owner);
            SetUpdateFunc(AIFunc.NormalUpdate);
        }

        public void SetUpdateFunc(UpdateDelegate handler)
        {
            _stateCtrl.UpdateHandler = handler;
        }

        public void Update()
        {
            //_stateCtrl.Update();
        }

        public void SetEnemy(Character enemy)
        {
            _enemyId = enemy.id;
            _owner.SetEnemy(enemy.id);
        }
    }
}// namespace Phoenix
