using System;

namespace Phoenix.Game.FightEmulator
{
    public static class MathUtil
    {   
        public static float RandomF(float min, float max)
        {
            return UnityEngine.Random.Range(min, max);
        }

        public static bool HitChance(float chance)
        {
            return RandomF(0, 1) <= chance;
        }

        public static int RandomI(int min, int max)
        {
            return UnityEngine.Random.Range(min, max+1);
        }
    }
} // Phoenix.Utils