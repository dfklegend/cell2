using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{	    
    // 创建战斗角色，方便调试战斗
    // 用于测试
    public class SimpleAttackCtrl : Singleton<SimpleAttackCtrl>
    {
        private Character _warrior;        
        public Character warrior { get { return _warrior; } }
        private Character _rogue;
        public Character rogue { get { return _rogue; } }
        //
        public void Prepare()
        {
            if (_warrior != null)
                return;
            Character warrior = new Character();
            warrior.SetName("战士");
            warrior.SetClass((int)eClass.Warrior);
            warrior.SetLevel(100);
            warrior.PrepareAttrs();
            warrior.attrs.Dump();

            Character rogue = new Character();
            rogue.SetName("盗贼");
            rogue.SetClass((int)eClass.Rogue);
            rogue.SetLevel(10);
            rogue.PrepareAttrs();
            rogue.attrs.Dump();

            _warrior = warrior;
            _rogue = rogue;
        }

        public void Attack()
        {
            FormulaResult result = FormulaUtil.SimpleFormula(null, _warrior, _rogue) as FormulaResult;           
            Debug.Log($"{_warrior.name}攻击{_rogue.name}，造成{(int)result.data.Dmg}点伤害");
        }
    }
    
} // namespace Phoenix
