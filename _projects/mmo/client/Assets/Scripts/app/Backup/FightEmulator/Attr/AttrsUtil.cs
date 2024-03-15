using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System.Text;

namespace Phoenix.Game.FightEmulator
{
    public static class AttrsUtil
    {   
        public static void CloneElement(AttrElement src, AttrElement tar)
        {
            tar.baseValue = src.baseValue;
            tar.percent = src.percent;
        }

        public static void CloneAttr(Attr src, Attr tar)
        {
            if (src == null || tar == null)
                return;
            CloneElement(src.GetBase(), tar.GetBase());
            CloneElement(src.GetAppend(), tar.GetAppend());
        }
        
        public static void Clone(Attrs src, Attrs tar)
        {
            for(var i = 0; i < AttrDefine.attrNames.Length; i ++)
            {
                var name = AttrDefine.attrNames[i];
                CloneAttr(src.GetAttr(name), tar.GetAttr(name));
            }
        }

        public static void OffsetElement(AttrElement main, float factor, AttrElement off)
        {
            main.baseValue += factor*off.baseValue;
            main.percent += factor*(off.percent - 1.0f);           
        }

        // main += factor*off
        public static void OffsetAttr(Attr main, float factor, Attr off)
        {
            if (main == null || off == null)
                return;
            OffsetElement(main.GetBase(), factor, off.GetBase());
            OffsetElement(main.GetAppend(), factor, off.GetAppend());
        }

        // main = main + factor*off
        public static void Offset(Attrs main, float factor,Attrs off)
        {
            for (var i = 0; i < AttrDefine.attrNames.Length; i++)
            {
                var name = AttrDefine.attrNames[i];
                OffsetAttr(main.GetAttr(name), factor, off.GetAttr(name));                
            }
        }

        public static eElementType GetElementType(string str)
        {            
            if (str == "a")
                return eElementType.Append;
            if (str == "all")
                return eElementType.All;
            return eElementType.Base;
        }

        public static (string attrName,bool percent, eElementType eleType) 
            AnalysisTypeString(string attrTypeString)
        {
            var subs = attrTypeString.Split(',');

            string attrName = "";
            bool percent = false;
            eElementType eType = eElementType.Base;
            if (subs.Length >= 3)
                eType = GetElementType(subs[2]);
            if (subs.Length >= 2)
                percent = subs[1] == "%";
            if (subs.Length >= 1)
                attrName = subs[0];
            return (attrName, percent, eType);
        }

        // 根据typeString
        // 力量,%,b
        public static void AddAttr(IAttrOwner owner, string attrTypeString,
             float value)
        {
            string attrName;
            bool percent;
            eElementType eType;
            (attrName, percent, eType) = AnalysisTypeString(attrTypeString);
            AddAttr(owner, attrName, percent, eType, value);
        }

        private static void addAttrWithNonPercentEType(IAttrOwner owner, string attrTypeString, float value,
            eElementType nonPercentEType)
        {
            string attrName;
            bool percent;
            eElementType eType;
            (attrName, percent, eType) = AnalysisTypeString(attrTypeString);

            if (!percent)
                eType = nonPercentEType;
            AddAttr(owner, attrName, percent, eType, value);
        }

        // 装备技能都是加base
        // buf才可能加append
        public static void AddEquipAttr(IAttrOwner owner, string attrTypeString, float value)
        {
            addAttrWithNonPercentEType(owner, attrTypeString, value, eElementType.Base);
        }

        public static void AddAttr(IAttrOwner owner, string attrName,
            bool percent, eElementType eleType, float value)
        {
            owner.AddAttr(attrName, percent, eleType, value);
        }

        public static void MakeAttrsBrief(StringBuilder sb, Attrs attrs)
        {
            briefAppendAttr(sb, attrs, AttrDefine.Strength);
            briefAppendAttr(sb, attrs, AttrDefine.Agility);
            briefAppendAttr(sb, attrs, AttrDefine.Intellect);
            briefAppendAttr(sb, attrs, AttrDefine.Stamina);

            briefAppendAttr(sb, attrs, AttrDefine.HPMax);
            
            briefAppendAttr(sb, attrs, AttrDefine.MeleePower);
            briefAppendAttr(sb, attrs, AttrDefine.PhysicHit);
            briefAppendAttr(sb, attrs, AttrDefine.PhysicCritical);
        }

        private static void briefAppendAttr(StringBuilder sb, 
            Attrs attrs, string attrName)
        {
            sb.AppendFormat("{0}: {1}\n", attrName, attrs.GetAttr(attrName).intFinal);
        }
    }
}// namespace Phoenix
