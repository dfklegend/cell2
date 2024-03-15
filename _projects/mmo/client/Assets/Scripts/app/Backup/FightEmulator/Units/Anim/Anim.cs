using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using UnityEngine.UI;
using Phoenix.Utils;

namespace Phoenix.Game.FightEmulator
{	
    public class PendingAnim
    {
        public string anim;
        public float time;
    }

    // 只有idle状态
    // 动作播放完毕就切换一次到idle
    public class AnimCtrl
    {
        private Animator _animator;
        private bool _speedDirt = false;
        private float _speed = 1f;

        private float _animEndTime = -1f;
        
        private string _curAnim;
        private float _timeStartAnim;
        private float _curAnimTime;
        public List<PendingAnim> _pendings = new List<PendingAnim>();
        
        public void Init(Animator animator)
        {
            _animator = animator;
        }

        public void Update()
        {
            updatePending();
            if (_animEndTime > 0 && !isAnimPlaying() && _pendings.Count == 0)
            {
                analysisSpeed();                
                _animEndTime = -1;
                _animator.CrossFade("idle", 0.1f);
            }
        }

        private void setSpeed(float f)
        {
            _animator.speed = f;
            _speed = f;
            Log.LogCenter.Default.Debug("setSpeed:{0}", f);
        }

        private bool isAnimPlaying()
        {
            return TimeUtil.Now() < _animEndTime;
        }

        public void PlayAnim(string anim, float time)
        {
            if(isAnimPlaying())
            {
                var one = new PendingAnim();
                one.anim = anim;
                one.time = time;
                _pendings.Add(one);
                _speedDirt = true;

                Log.LogCenter.Default.Debug("push: {0} count:{1}", anim, _pendings.Count);
                analysisSpeed();
                return;
            }
            _doAnim(anim, time);
        }

        private void _doAnim(string anim, float time)
        {
            Log.LogCenter.Default.Debug("_doAnim: {0}", anim);
            // 如果发现之前动作没结束，清除一下
            var state = _animator.GetCurrentAnimatorStateInfo(0);
            if (state.IsName(anim))
            {
                // clear cur state
                clearAnimState();
            }
            _animator.CrossFade(anim, 0.1f);
            _timeStartAnim = TimeUtil.Now();
            _curAnimTime = time;
            _animEndTime = _timeStartAnim + time/_speed;
        }        

        private void clearAnimState()
        {
            _animator.Play("");            
        }

        private float calcSpeed()
        {
            return 1f + (_pendings.Count) * 0.5f;            
        }

        private void analysisSpeed()
        {
            if (!_speedDirt)
                return;
            _speedDirt = false;
            setSpeed(calcSpeed());
            // 重置下当前时间
            if(isAnimPlaying())
                _animEndTime = _timeStartAnim + _curAnimTime/_speed;
        }

        private void updatePending()
        {
            if (_pendings.Count == 0 || isAnimPlaying())
                return;
            var one = _pendings[0];
            _pendings.RemoveAt(0);

            _speedDirt = true;
            analysisSpeed();
            _doAnim(one.anim, one.time);
        }
    }    
} // namespace Phoenix
