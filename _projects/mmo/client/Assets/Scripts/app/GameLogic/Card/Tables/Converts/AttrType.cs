using Phoenix.csv;
using System;
using System.Collections.Generic;
using UnityEngine;

// 定义一些非常规数据的转换
namespace Phoenix.Game.Card
{
    public class AttrTypeConvert : BaseValueConvert
    {   
        public static AttrType defaultValue = new AttrType() {
            type = eAttrType.Invalid,
        };

        public override Type GetTargetType()
        {
            return typeof(AttrType);
        }
        public override object Convert(string value)
        {
            try
            {
                value = value.Trim();
                if (string.IsNullOrEmpty(value))
                    return defaultValue;
                var attr = new AttrType();
                string[] subs = value.Split(',');

                if(subs.Length >= 2)
                {
                    attr.percent = subs[1] == "%";
                }
                var type = StringsUtils.GetIndexFromArray(Strings.ATTR_NAMES, subs[0]);                
                attr.type = (eAttrType)type;
                return attr;
            }
            catch (System.Exception)
            {
                return null;
            }
        }
    }    
}