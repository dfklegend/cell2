using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    public class FormulaFactory : Singleton<FormulaFactory>
    {
        private FormulaSimpleTest _simple = new FormulaSimpleTest();
        private WeaponMelee _mainWeaponMelee = new WeaponMelee();


        private Dictionary<eFormulaType, IFormula> _formulas
            = new Dictionary<eFormulaType, IFormula>();
        private IFormula _default;

        FormulaFactory()
        {
            initFormulas();
        }

        private void initFormulas()
        {
            _formulas[eFormulaType.SimpleTest] = new FormulaSimpleTest();
            _formulas[eFormulaType.NoDmg] = new FormulaNoDmg();
            _formulas[eFormulaType.PhysicMelee] = new WeaponMelee();


            _default = _formulas[eFormulaType.SimpleTest];
        }

        public IFormula GetFormula(int formulaType)
        {
            IFormula ret;
            if (_formulas.TryGetValue((eFormulaType)formulaType, out ret))
                return ret;
            return _default;
        }
    }
}// namespace Phoenix
