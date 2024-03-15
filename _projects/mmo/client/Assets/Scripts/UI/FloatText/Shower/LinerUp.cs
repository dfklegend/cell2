using System.Collections.Generic;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game.FText
{
    [IntType(ShowType.LinerUp)]
    public class LinerUp : FloatTextShower
    {
        private float _life = 1f;
        private float _speed = 1f;
        private float _deathTime = 0f;
        private float _lastTime = 0f;
        // life, speed
        public void Create(params object[] args)
        {
            if (args.Length < 2)
                return;
            _life = (float)args[0];
            _speed = (float)args[1];

            _deathTime = Time.time + _life;
            _lastTime = Time.time;
        }

        public void Update(FloatText text)
        {
            float off = Time.time - _lastTime;
            if (off <= 0)
                return;
            _lastTime = Time.time;

            Vector2 pos = text.GetPos();
            pos.y += _speed * off;

            text.SetPos(pos);
        }

        public bool IsOver()
        {
            return Time.time >= _deathTime;
        }
    }
} // namespace Phoenix
