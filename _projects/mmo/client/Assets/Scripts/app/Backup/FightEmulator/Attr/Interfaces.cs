using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 定义属性系统的接口
    // 属性所有者
    // 比如: 角色
    public interface IAttrOwner
    {
        int GetLevel();
        // 获取属性对象
        Attrs GetAttrs();
        // 重新计算属性转化
        void RecalcAttrTransform();
        // 所有实际属性变化都通过这个接口添加
        // 装备,buf,其他
        void AddAttr(string attrName,
            bool percent, eElementType eleType, float value);
        void AddTransformAttr(string attrName, float value);

        string GetMonsterId();
        IWeapon GetEnemyWeapon();
    }

    // 抽象owner中需要访问的外部信息
    public interface IOwnerImpl
    {
        int GetLevel();
        void OnAttrChanged(string attrName, float oldV, float newV);
    }

    // 属性特质
    // 比如: 种族，职业
    // 假设可以兼职，传入对应职业等级即可
    public interface IAttrFeature
    {
        // 初始化才会调用
        void ApplyInitAttrs(IAttrOwner owner);
        // 初始化和升级时调用
        void ApplyLevelupAttrs(IAttrOwner owner, int offset);
    }

    
    // 可以动态插拔的属性
    // 比如 装备，buf，天赋符文等
    public interface IAttrEquipable
    {
        void Equip(IAttrOwner owner);
        void Unequip(IAttrOwner owner);
    }

    // 属性转化器
    public interface IAttrTransformer
    {
        string GetAttrName();
        // 属于几阶
        int GetLevel();
        // 转化
        void Transform(IAttrOwner owner);
    }
}// namespace Phoenix
