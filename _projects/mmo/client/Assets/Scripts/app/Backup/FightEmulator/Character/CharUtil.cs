using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System.Text;

namespace Phoenix.Game.FightEmulator
{
    public static class CharUtil
    {
        public static bool CanCharMoving(Character owner)
        {
            if (owner.IsDead())
                return false;
            if (owner.HasMainSkillRunning())
                return false;
            return true;
        }

        public static float DistTo(Character src, Character tar)
        {
            Vector3 off = src.pos - tar.pos;
            return off.magnitude;
        }

        public static Vector3 GetMoveTar(Character src, Character tar, float tarDist)
        {
            var dir = tar.pos - src.pos;
            dir.Normalize();
            return src.pos + dir * tarDist;
        }
    }
}// namespace Phoenix
