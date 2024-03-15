using Phoenix.csv;
using System;
using System.Collections.Generic;
using UnityEngine;

// 定义一些非常规数据的转换
namespace Phoenix.Game.FightEmulator
{
    public class StringArrayConvert : BaseValueConvert
    {   
        public override Type GetTargetType()
        {
            return typeof(string[]);
        }
        public override object Convert(string value)
        {
            try
            {
                value = value.Trim();
                if (string.IsNullOrEmpty(value))
                    return null;
                string[] subs = value.Split(',');
                if (subs.Length == 1 && string.IsNullOrEmpty(subs[0]))
                    return null;
                return subs;
            }
            catch (System.Exception)
            {
                return null;
            }
        }
    }    
}