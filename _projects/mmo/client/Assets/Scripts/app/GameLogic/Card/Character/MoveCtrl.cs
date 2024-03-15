using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System.Text;

namespace Phoenix.Game.Card
{
    public class MoveCtrl
    {
        Character _owner;
        float _speed = 2.0f;

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

        public float calcViewSpeed(bool moveAlignX)
        {
            // 由于高度高一些，所以速度缩减一点，避免
            if(moveAlignX)
                return _speed;
            return _speed * Consts.GridSizeWidth / Consts.GridSizeHeight;
        }

        public void Update()
        {
            if (!_moving)
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
