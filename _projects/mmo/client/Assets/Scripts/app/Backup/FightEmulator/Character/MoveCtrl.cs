using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System.Text;

namespace Phoenix.Game.FightEmulator
{
    public class MoveCtrl
    {
        Character _owner;
        float _speed = 1.0f;

        bool _moving = false;
        Vector3 _tar = Vector3.zero;

        public void Init(Character owner)
        {
            _owner = owner;
        }
        
        public void MoveTo(float x, float z)
        {
            _moving = true;
            _tar.x = x;
            _tar.z = z;
        }

        public void Stop()
        {
            _moving = false;
        }

        public void Update()
        {
            if (!_moving)
                return;
            if (!CharUtil.CanCharMoving(_owner))
                return;
            float delta = Time.deltaTime;

            float step = delta * _speed;

            Vector3 dir = _tar - _owner.pos;
            var dist = dir.magnitude;

            if(step >= dist)
            {
                SetPos(_tar);
                _moving = false;
                return;
            }

            dir.Normalize();
            Vector3 newPos = _owner.pos + dir * step;
            SetPos(newPos);
        }

        private void SetPos(Vector3 newPos)
        {
            _owner.pos = newPos;
        }
    }
}// namespace Phoenix
