using Phoenix.csv;
using System;
using System.Collections.Generic;
using UnityEngine;

// 定义一些非常规数据的转换
namespace Phoenix.Game.FightEmulator
{
    public class FormulaConvert : BaseValueConvert
    {
        private StringToInt _map = new StringToInt();

        public FormulaConvert()
        {
            _map.InitFromStrings(FormulaDefine.formulaTypes, 0);
        }

        public override Type GetTargetType()
        {
            return typeof(eFormulaType);
        }
        public override object Convert(string value)
        {
            try
            {
                return _map.ToInt(value);
            }
            catch (System.Exception)
            {
                return null;
            }
        }
    }    
}