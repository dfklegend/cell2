using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    public class FormulaNoDmg : IFormula
    {
        // 无伤害
        public IFormulaResult Apply(ISkill skill, ICharacter srcChar, ICharacter tarChar)
        {
            return null;
        }
    }
}// namespace Phoenix
