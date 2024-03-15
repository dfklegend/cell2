using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

using Phoenix.Game.FightEmulator;


namespace Phoenix.Game.Skill
{	    
    public enum eHandType
    {
        None = -1,
        MainHand = 0,
        OffHand,
        Max
    }

    public class HandData
    {
        // -1代表无效
        public float interval;
        public float nextCanAttack;
    }

    // 攻击计时
    // 主副手
    public class PhysicAttackCtrl
    {
        private HandData[] _hands = new HandData[(int)eHandType.Max];

        public PhysicAttackCtrl()
        {
            for(int i = 0; i < (int)eHandType.Max; i ++)
            {
                _hands[i] = new HandData();
                _hands[i].interval = -1;
            }
        }
        
        public void SetHandAttackInterval(eHandType handType, float v)
        {
            var hand = getHand(handType);
            if (hand == null)
                return;
            hand.interval = v;
        }

        public void SetHandNoWeapon(eHandType handType)
        {
            SetHandAttackInterval(handType, - 1);
        }

        private HandData getHand(eHandType hand)
        {
            if (hand < eHandType.MainHand || hand >= eHandType.Max)
                return null;
            return _hands[(int)hand];
        }

        public bool CanAttack(eHandType handType, float now)
        {
            var hand = getHand(handType);
            if (hand == null)
                return false;
            if (hand.interval < 0)
                return false;
            return now >= hand.nextCanAttack;
        }

        public void OnAttackStart(eHandType handType, float now)
        {
            var hand = getHand(handType);
            if (hand == null)
                return;
            hand.nextCanAttack = now + hand.interval;
        }        

        // 被打断
        public void OnAttackStop(eHandType handType)
        {
            var hand = getHand(handType);
            if (hand == null)
                return;
            hand.nextCanAttack = 0;
        }

        // 避免主副手一起触发
        public void SureAttackDisplayDelay(eHandType handType, float now)
        {
            const float AttackDisplayDelay = 0.2f;
            var hand = getHand(handType);
            if (hand == null)
                return;            
            if (hand.interval < 0)
                return;
            // 避免0.2s内连续触发主副手
            if(now >= hand.nextCanAttack + AttackDisplayDelay)
            {
                hand.nextCanAttack = now + AttackDisplayDelay;
            }
        }
    }
} // namespace Phoenix
