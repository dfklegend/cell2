using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using System.Net.Sockets;
using Network;


namespace Phoenix.Game.Card
{   
    public class HEventTestSnapshot : HEvent<int>
    {
        public Cproto.TestSnapshot snapshot;
        public HEventTestSnapshot(Cproto.TestSnapshot snapshot)
            : base(EventDefine.TestSnapshot)
        {
            this.snapshot = snapshot;
        }
    }    

    public class HEventMonsterSnapshot : HEventServerMsg<Cproto.MonsterSnapshot>
    {
        public HEventMonsterSnapshot(Cproto.MonsterSnapshot msg)
            : base(msg, EventDefine.MonsterSnapshot)
        {
        }
    }

    public class HEventExitSnapshot : HEventServerMsg<Cproto.ExitSnapshot>
    {
        public HEventExitSnapshot(Cproto.ExitSnapshot msg)
            : base(msg, EventDefine.ExitSnapshot)
        {
        }
    }

    public class HEventMoveTo : HEvent<int>
    {
        public Cproto.MoveTo msg;
        public HEventMoveTo(Cproto.MoveTo msg)
            : base(EventDefine.MoveTo)
        {
            this.msg = msg;
        }
    }

    public class HEventStopMove : HEvent<int>
    {
        public Cproto.MoveTo msg;
        public HEventStopMove(Cproto.MoveTo msg)
            : base(EventDefine.StopMove)
        {
            this.msg = msg;
        }
    }

    public class HEventUnitLeave : HEvent<int>
    {
        public Cproto.UnitLeave msg;
        public HEventUnitLeave(Cproto.UnitLeave msg)
            : base(EventDefine.UnitLeave)
        {
            this.msg = msg;
        }
    }

    //public class HEventAttack : HEvent<int>
    //{
    //    public Cproto.Attack msg;
    //    public HEventAttack(Cproto.Attack msg)
    //        : base(EventDefine.Attack)
    //    {
    //        this.msg = msg;
    //    }
    //}

    public class HEventRelive : HEvent<int>
    {
        public Cproto.UnitRelive msg;
        public HEventRelive(Cproto.UnitRelive msg)
            : base(EventDefine.UnitRelive)
        {
            this.msg = msg;
        }
    }    

    // 服务器开始技能
    public class HEventServerStartSkill : HEventServerMsg<Cproto.StartSklill>
    {        
        public HEventServerStartSkill(Cproto.StartSklill msg)
            : base(msg, EventDefine.ServerStartSkill)
        {            
        }
    }

    // 服务器技能命中
    public class HEventServerSkillHit : HEventServerMsg<Cproto.SkillHit>
    {
        public HEventServerSkillHit(Cproto.SkillHit msg)
            : base(msg, EventDefine.ServerSkillHit)
        {
        }
    }

    public class HEventServerSkillBroken : HEventServerMsg<Cproto.SkillBroken>
    {
        public HEventServerSkillBroken(Cproto.SkillBroken msg)
            : base(msg, EventDefine.ServerSkillBroken)
        {
        }
    }

    // 表现
    public class HEventStartSkill : HEvent<int>
    {
        public Character src;
        public string skillId;
        public int tarId;


        public HEventStartSkill(Character src, string skillId, int tarId)
            : base(EventDefine.StartSkill)
        {
            this.src = src;
            this.skillId = skillId;
            this.tarId = tarId;
        }
    }

    // 同步血条
    public class HEventHPChanged : HEvent<int>
    {
        public Character src;

        public HEventHPChanged(Character src)
            : base(EventDefine.HPChanged)
        {
            this.src = src;
        }
    }

    public class HEventMPChanged : HEvent<int>
    {
        public Character src;

        public HEventMPChanged(Character src)
            : base(EventDefine.MPChanged)
        {
            this.src = src;
        }
    }

    // 表现
    public class HEventAttack : HEvent<int>
    {
        public string skillId;
        public Character src;
        public Character tar;
        public FormulaResult result;


        public HEventAttack(string skillId, Character src, Character tar, FormulaResult result)
            : base(EventDefine.Attack)
        {
            this.skillId = skillId;
            this.src = src;
            this.tar = tar;
            this.result = result;
        }
    }

    public class HEventSkillBroken : HEvent<int>
    {
        public Character src;
        public string skillId;    


        public HEventSkillBroken(Character src, string skillId)
            : base(EventDefine.SkillBroken)
        {
            this.src = src;
            this.skillId = skillId;           
        }
    }

}// namespace Phoenix
