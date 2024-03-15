using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{   
    public class ClassTransformerBuilder : Singleton<ClassTransformerBuilder>
    {
        public void BuildTransformer(CharAttrs charAttrs, int classType)
        {
            // common
            charAttrs.RegisterAttrTransformer(HPMaxTransformer.It);
            charAttrs.RegisterAttrTransformer(ArmorTransformer.It);
            // class feature
            switch (classType)
            {
                case (int)eClass.Warrior:
                    charAttrs.RegisterAttrTransformer(WarriorMeleePowerTransformer.It);
                    break;
                case (int)eClass.Rogue:
                    charAttrs.RegisterAttrTransformer(RogueMeleePowerTransformer.It);
                    break;
            }
        }
    }
}// namespace Phoenix
