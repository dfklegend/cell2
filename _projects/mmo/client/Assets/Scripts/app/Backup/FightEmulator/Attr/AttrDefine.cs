using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    public static class AttrDefine
    {
        // 使用名字，便于扩展脚本使用
        // 下面定义便于程序使用
        public static string Strength = "力量";
        public static string Agility = "敏捷";
        public static string Intellect = "智力";
        public static string Stamina = "耐力";
        public static string HP = "HP";
        public static string HPMax = "HPMax";

        public static string MeleePower = "近战攻强";
        public static string RangePower = "远程攻强";
        public static string SpellPower = "法强";

        public static string PhysicHit = "物理命中";
        public static string PhysicCritical = "物暴点数";
        public static string SpellHit = "法术命中";
        public static string SpellCritical = "法暴点数";
        public static string Armor = "护甲";
        public static string SpellArmor = "法术护甲";
        // 
        public static string ShieldBlock = "盾牌格挡";
        public static string BlockChance = "格挡点数";
        public static string WeaponEnhance = "武器增强";
        

        // 力量*2转化为攻击
        // 力量*1转化为防御
        private static string[] _attrNames = {
            Strength, Agility, Intellect, Stamina,
            HP, HPMax,
            MeleePower, RangePower, SpellPower,
            PhysicHit, PhysicCritical, SpellHit, SpellCritical,
            Armor, SpellArmor, ShieldBlock, BlockChance, WeaponEnhance,            
        };

        public static string[] attrNames { get { return _attrNames; } }

        public static void InitAttrs(Attrs attrs)
        {            
            for(var i = 0; i < _attrNames.Length; i ++)
            {
                attrs.NewAttr(_attrNames[i]);
            }
        }      
    }

    public enum eElementType
    {
        Base,           // 
        Append,
        //Transformed,    // 转化
        All,
    }
        
}// namespace Phoenix
