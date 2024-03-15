using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

namespace Phoenix.Game.FightEmulator
{	   
    // 向目标飞行子弹
    public class Bullet
    {
        private string _skillId;
        private float _speed = 2f;
        private int _tarId = -1;
        private int _ownerId = -1;
        private bool _over = false;
        private Vector3 _pos;
        private int _id;
        public int id { get { return _id; } }

        public Vector3 pos { get { return _pos; } }

        public Bullet(int id, int ownerId, int tarId, string skillId, float speed)
        {
            _id = id;
            _ownerId = ownerId;
            _tarId = tarId;
            _skillId = skillId;
            _speed = speed;
            var owner = FightCtrl.It.GetChar(_ownerId);
            _pos = owner.pos;
        }

        public bool IsOver() 
        {
            return _over;
        }

        private void setOver()
        {
            _over = true;
        }

        public void Update()
        {
            if (IsOver())
                return;
            Character owner;
            Character tar;
            if((owner = FightCtrl.It.GetChar(_ownerId)) == null ||
                (tar = FightCtrl.It.GetChar(_tarId)) == null)
            {
                setOver();
                return;
            }

            float delta = Time.deltaTime;

            float step = _speed * delta;
            var offset = tar.pos - _pos;
            var dist = offset.magnitude;

            bool hitted = false;
            if(step >= dist)
            {
                step = dist;
                hitted = true;
            }

            var dir = offset;
            dir.Normalize();

            _pos = _pos + dir * step;
            if(hitted)
            {
                doHit();
                setOver();
            }
        }

        private void doHit()
        {
            Log.LogCenter.Default.Debug("{0} skill {1} bullet hit {2}",
                _ownerId, _skillId, _tarId);

            var owner = FightCtrl.It.GetChar(_ownerId);
            if (owner == null)
                return;
            var ext = new Skill.SkillArgsEx();
            ext.onlyBalance = true;
            owner.skillCtrl.BGSkillToTar(_tarId, _skillId, ext);
        }
    }
} // namespace Phoenix
