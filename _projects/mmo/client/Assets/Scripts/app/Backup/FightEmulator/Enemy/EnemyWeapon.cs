using UnityEngine;
using System.Reflection;
using Phoenix.csv;

namespace Phoenix.Game.FightEmulator
{
    public class EnemyWeapon : IWeapon
    {
        private int _minDmg = 0;
        private int _maxDmg = 0;
        private float _speed = 3f;
        
        public void OffMinDmg(float off)
        {
            _minDmg += (int)off;
        }
               
        public int GetMinDmg()
        {
            return _minDmg;
        }

        public void OffMaxDmg(float off)
        {
            _maxDmg += (int)off;
        }

        public int GetMaxDmg()
        {
            return _maxDmg;
        }

        public float GetSpeed()
        {
            return _speed;
        }

        public void SetSpeed(float v)
        {
            _speed = v;
        }

        public void Reset()
        {
            _minDmg = 0;
            _maxDmg = 0;
        }
    }
}
