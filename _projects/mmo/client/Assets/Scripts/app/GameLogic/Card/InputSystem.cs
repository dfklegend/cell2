using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.Card
{ 
    // 简单的操控
    public class InputSystem : Singleton<InputSystem>
    {
        private Character _player;
        
        private bool _moving;
        private Vector3 _dir;
        private Vector3 _lastTar;
        private float _lastSend;

        public void Reset()
        {
            _dir = Vector3.zero;
            _lastTar = Vector3.zero;
            _moving = false;
        }

        public void OnPlayerEnter(Character player)
        {
            _player = player;
            Reset();
        }

        public void OnPlayerLeave()
        {
            _player = null;
        }

        public void Update()
        {
            if (_player == null)
                return;
            Vector3 off = Vector3.zero;
            if(Input.GetKey(KeyCode.LeftArrow))
            {
                off.x = -1;
            }
            if (Input.GetKey(KeyCode.RightArrow))
            {
                off.x = 1;
            }
            if (Input.GetKey(KeyCode.UpArrow))
            {
                off.z = 1;
            }
            if (Input.GetKey(KeyCode.DownArrow))
            {
                off.z = -1;
            }

            off.Normalize();
            var moving = off.sqrMagnitude > 0.1;
            if(_moving && !moving)
            {
                onStopMove();
            }
            if(!_moving && moving)
            {
                onStartMove(off);
            }
            _moving = moving;
            if (moving)
            {                
                updateMoving(off);
            }
        }

        private void onStartMove(Vector3 dir)
        {
            Debug.Log("onStartMove: " + dir);
        }

        private void onStopMove()
        {
            Debug.Log("onStopMove");
            _dir = Vector3.zero;
            sendStopMove();
        }

        private void updateMoving(Vector3 dir)
        {
            // 如果方向变更，立刻发送新请求
            // 否则每隔一段时间更新新目标(离上次请求目标比较近的时候)
            if(_dir != dir)
            {
                Debug.Log("onMove: " + dir);
                _dir = dir;
                updateTarAndSend(dir);
                return;
            }

            var newTar = _player.pos + dir * 3.0f;            
            var offset = _lastTar - newTar;
            // 目标改变了
            if(offset.sqrMagnitude > 0.5f || Time.time > _lastSend + 1.0)
            {
                updateTarAndSend(dir);
            }
        }

        private void updateTarAndSend(Vector3 dir)
        {
            var tar = _player.pos + dir * 3.0f;
            _lastTar = tar;
            _lastSend = Time.time;
            Systems.It.GetSystem<ControlSystem>("control").Request<Cproto.EmptyArg>("moveto",
                new Cproto.ReqMoveTo {X = tar.x, Z = tar.z }, null);
        }

        private void sendStopMove()
        {
            var tar = _player.pos;
            Systems.It.GetSystem<ControlSystem>("control").Request<Cproto.EmptyArg>("stopmove",
                new Cproto.ReqMoveTo { X = tar.x, Z = tar.z }, null);
        }
    }

} // namespace Phoenix
