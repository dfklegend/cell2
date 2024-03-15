using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 属性计算部分统一放这里
    // character集成使用
    public class CharAttrs : IAttrOwner
    {
        const int MaxAttrTransformLevel = 3;

        // 提供对外部属性的访问
        private IOwnerImpl _impl;
        private Attrs _attrs;
        public Attrs attrs { get { return _attrs; } }

        // 属性特性
        private List<IAttrFeature> _features = new List<IAttrFeature>();

        // 属性转化
        // attrName:
        private Dictionary<string, IAttrTransformer> _transformers =
            new Dictionary<string, IAttrTransformer>();

        // 一些快捷方式
        private Attr _hp;
        private Attr _hpMax;

        private string _monsterId;
        // 虚拟的敌人武器
        private IWeapon _enemyWeapon;

        public void InitAttrs(IOwnerImpl impl, Attrs attrs)
        {
            _impl = impl;
            _attrs = attrs;

            _hp = attrs.GetAttr(AttrDefine.HP);
            _hpMax = attrs.GetAttr(AttrDefine.HPMax);
        }

        public void InitEnemy(string monsterId)
        {
            _monsterId = monsterId;
            _enemyWeapon = new EnemyWeapon();
        }

        public int GetLevel()
        {
            return _impl.GetLevel();
        }

        public Attrs GetAttrs()
        {
            return _attrs;
        }

        public int GetHP()
        {
            return _hp.intFinal;
        }

        public void SetHP(int hp)
        {
            _hp.Base.baseValue = (float)hp;
        }

        public int GetHPMax()
        {
            return _hpMax.intFinal;
        }

        public float GetHPPercent()
        {
            if (_hpMax.intFinal <= 0)
                return 1;
            return (float)_hp.intFinal / _hpMax.intFinal;
        }

        public void RecalcAttrTransform()
        {
            // 清除高阶属性的transform value
            clearAllAttrTransformValue();
            // 计算转化的属性
            for (var i = 2; i <= MaxAttrTransformLevel; i++)
            {
                CalcAttrTransformByLevel(i);
            }
        }

        private void clearAllAttrTransformValue()
        {
            foreach (KeyValuePair<string, IAttrTransformer> kv in _transformers)
            {
                var attr = _attrs.GetAttr(kv.Key);
                if (attr == null)
                    continue;
                attr.Transformed.ResetValue();
            }
        }

        // 计算不同阶的数值转化
        private void CalcAttrTransformByLevel(int level)
        {            
            foreach(KeyValuePair<string, IAttrTransformer> kv in _transformers)
            {
                var transformer = kv.Value;
                if (transformer.GetLevel() != level)
                    continue;
                transformer.Transform(this);
            }
        }
        
        public void FullyCalcAttrs()
        {
            // 重新计算
            _attrs.Reset();

            // 每个feature
            initFeaturesAttr();            

            RecalcAttrTransform();
        }

        private void initFeaturesAttr()
        {
            foreach(var one in _features)
            {
                initFeatureAttr(one);
            }
        }

        // 如果出现职业兼职
        // 提供一个getFeatureLevel的抽象
        private void initFeatureAttr(IAttrFeature feature)
        {
            feature.ApplyInitAttrs(this);
            if (GetLevel() > 1)
                feature.ApplyLevelupAttrs(this, GetLevel() - 1);
        }

        public void AddFeature(IAttrFeature feature)
        {
            _features.Add(feature);
        }

        public void RegisterAttrTransformer(IAttrTransformer transformer)
        {
            _transformers[transformer.GetAttrName()] = transformer;
        }

        public void Levelup(int offset)
        {
            foreach (var one in _features)
            {
                one.ApplyLevelupAttrs(this, offset);
            }
            RecalcAttrTransform();
        }

        // 装备
        public void Equip(IAttrEquipable equipable)
        {
            equipable.Equip(this);
            RecalcAttrTransform();
        }

        public void Unequip(IAttrEquipable equipable)
        {
            equipable.Unequip(this);
            RecalcAttrTransform();
        }
        
        public void AddAttr(string attrName,
            bool percent, eElementType eleType, float value)
        {
            var attr = _attrs.GetAttr(attrName);
            if (attr == null)
                return;

            var oldV = attr.final;
            applyAddAttr(attr, percent, eleType, value);
            var newV = attr.final;

            // 数值变化
            _impl.OnAttrChanged(attrName, oldV, newV);
        }

        private void applyAddAttr(Attr attr,
            bool percent, eElementType eleType, float value)
        {            
            if (percent)
            {
                bool addBase = eleType == eElementType.Base ||
                    eleType == eElementType.All;
                bool addAppend = eleType == eElementType.Append ||
                    eleType == eElementType.All;
                if (addBase)
                    attr.Base.percent += value;
                if (addAppend)
                    attr.Append.percent += value;
                return;
            }

            // 基本属性, all无效
            if (eleType == eElementType.Base)
                attr.Base.baseValue += value;
            else
                attr.Append.baseValue += value;            
        }

        public void AddTransformAttr(string attrName, float value)
        {
            var attr = _attrs.GetAttr(attrName);
            if (attr == null)
                return;
            var oldV = attr.final;
            applyAddTransformAttr(attr, value);
            var newV = attr.final;

            // 数值变化
            _impl.OnAttrChanged(attrName, oldV, newV);
        }

        private void applyAddTransformAttr(Attr attr, float value)
        {            
            attr.Transformed.baseValue += value;
        }

        public IWeapon GetEnemyWeapon()
        {
            return _enemyWeapon;
        }

        public string GetMonsterId()
        {
            return _monsterId;
        }
    }
}// namespace Phoenix
