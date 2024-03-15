using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 结算中用的中间数据
    // 结算过程中，会抛出结算事件
    //     相关的机制 可以通过监听事件，对结算数值做修改
    // 比如，某些机制可能额外造成伤害
    // 额外提供护甲穿透
    public class FormulaOpData
    {
        public int srcId;
        public int tarId;

        public bool hit = false;
        public bool critical = false;
        public bool block = false;

        public Skill.eHandType hand;
        
        public float Dmg = 0;
        public int weaponDmg = 0;
        // 格挡掉的伤害
        public int dmgBlocked = 0;
    }

    // 结算结果对象    
    public class FormulaResult : IFormulaResult
    {
        public FormulaOpData data = new FormulaOpData();
    }
}// namespace Phoenix
