using System;
using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;


namespace Phoenix.csv
{
    // 可以定制数值转换类型
    public abstract class BaseValueConvert
    {
        public virtual Type GetTargetType() { return null; }
        public virtual object Convert(string value) { return null; }
    }

    public class IntValueConvert : BaseValueConvert
    {
        public override Type GetTargetType()
        {
            return typeof(int);
        }
        public override object Convert(string value)
        {
            try
            {
                return int.Parse(value);
            }
            catch (System.Exception)
            {
                return 0;
            }
        }
    }

    public class FloatValueConvert : BaseValueConvert
    {
        public override Type GetTargetType()
        {
            return typeof(float);
        }
        public override object Convert(string value) 
        {
            try
            {
                return float.Parse(value);
            }
            catch (System.Exception)
            {                
                return 0f;
            }
        }
    }

    public class DoubleValueConvert : BaseValueConvert
    {
        public override Type GetTargetType()
        {
            return typeof(double);
        }
        public override object Convert(string value)
        {
            try
            {
                return double.Parse(value);
            }
            catch (System.Exception)
            {                
                return 0f;
            }
        }
    }

    public class CSVValueConverter : Singleton<CSVValueConverter>
    {
        Dictionary<System.Type, BaseValueConvert> _convert = new Dictionary<Type, BaseValueConvert>();

        public CSVValueConverter()
        {
            InitAll();
        }

        void InitAll()
        {
            List<Type> cs = SystemUtil.GetAllClass<BaseValueConvert>();
            if (cs == null)
                return;
            for (int i = 0; i < cs.Count; i++)
            {
                BaseValueConvert c = System.Activator.CreateInstance(cs[i]) as BaseValueConvert;
                if( c == null )
                    continue;
                _convert[c.GetTargetType()] = c;
            }
        }

        public BaseValueConvert FindConvert( System.Type t )
        {
            BaseValueConvert c;
            if (_convert.TryGetValue(t, out c))
                return c;
            return null;
        }
    }
}