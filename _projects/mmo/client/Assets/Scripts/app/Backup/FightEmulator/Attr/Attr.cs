using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;

namespace Phoenix.Game.FightEmulator
{	
    // 为了保证属性增加减少的无损，实际计算中不限制
    // 有效性范围，只是在计算最终的final时，提供合法性处理
    // 可以用外部定制属性的有效范围
    // 比如: 生命衰减最少到多少
    public interface IElementValueClamper
    {
        float ClampBase(float v);
        float ClampPercent(float v);
        float ClampFinal(float v);
    }

    // 属性最终值
    // 比如HPMax最小为1
    public interface IAttrFinalClamper
    {
        float ClampFinal(float v);
    }


    public class AttrElement
    {
        private IElementValueClamper _clamper;
        float _base = 0.0f;
        float _percent = 1.0f;
        float _final = 0.0f;
        bool _dirt = false;

        public void SetClamper(IElementValueClamper c)
        {
            _clamper = c;
        }

        public float final 
        { 
            get
            {
                tryCalcFinal();
                return _final;
            }
        }

        public bool dirt
        {
            get { return _dirt; }
        }

        public float baseValue 
        {
            set { _base = value; _dirt = true; }
            get { return _base; }
        }

        public float percent
        {
            set { _percent = value; _dirt = true; }
            get { return _percent; }
        }

        private void tryCalcFinal()
        {
            if (!_dirt)
                return;
            _dirt = false;          
            // base * percent
            _final = clampFinalValue(clampBaseValue() * clampPercentValue());
        }

        private float clampBaseValue()
        {
            if (_clamper != null)
                return _clamper.ClampBase(_base);
            return Math.Max(_base, 0.0f);
        }

        private float clampPercentValue()
        {
            if (_clamper != null)
                return _clamper.ClampPercent(_percent);            
            // 逻辑意义上最少减到10%属性值
            return Math.Max(_percent, 0.1f);
        }

        private float clampFinalValue(float v)
        {
            if (_clamper != null)
                return _clamper.ClampFinal(v);
            return v;
        }

        public void Reset()
        {
            _base = 0;
            _percent = 1;
            _final = 0;
            _dirt = false;
        }

        public void ResetValue()
        {
            _base = 0;            
            _final = 0;
            _dirt = false;
        }
    }

    public class AttrConst
    {
        // 角色基础属性
        // 来自初始，成长，装备
        public const int ElementBase = 0;
        // 附加的数值
        // 来自buf等
        public const int ElementAppend = 1;
        // 由其他属性转化的属性
        // 比如 力量转化为攻强
        public const int ElementTransformed = 2;
        public const int MaxAttrElement = 3;
    }    

    public class Attr
    {
        private IAttrFinalClamper _clamper;
        private AttrElement[] _elements = null;
        private float _final = 0.0f;

        public Attr()
        {
            _elements = new AttrElement[AttrConst.MaxAttrElement];
            for (var i = 0; i < AttrConst.MaxAttrElement; i++)
                _elements[i] = new AttrElement();
        }

        public void SetClamper(IAttrFinalClamper c)
        {
            _clamper = c;
        }

        public AttrElement GetBase()
        {
            return _elements[AttrConst.ElementBase];
        }

        public AttrElement Base
        {
            get { return _elements[AttrConst.ElementBase];  }
        }

        public AttrElement GetAppend()
        {
            return _elements[AttrConst.ElementAppend];
        }

        public AttrElement Append
        {
            get { return _elements[AttrConst.ElementAppend]; }
        }

        public AttrElement GetTransformed()
        {
            return _elements[AttrConst.ElementTransformed];
        }

        public AttrElement Transformed
        {
            get { return _elements[AttrConst.ElementTransformed]; }
        }

        public float final
        {
            get 
            {
                if (dirt)
                {
                    _final = calcFinal();
                }
                return _final;
            }            
        }

        private bool dirt
        {
            get { return Base.dirt || Append.dirt || Transformed.dirt; }
        }

        private float calcFinal()
        {
            float v = Base.final + Append.final + Transformed.final;
            if (_clamper != null)
                return _clamper.ClampFinal(v);
            return v;
        }

        public int intFinal
        {
            get { return (int)final; }
        }        

        public void Reset()
        {
            for (var i = 0; i < _elements.Length; i++)
                _elements[i].Reset();
        }

        public void VisitElements(Action<AttrElement> filter)
        {
            for (var i = 0; i < _elements.Length; i++)
                filter(_elements[i]);
        }
    }
}// namespace Phoenix
