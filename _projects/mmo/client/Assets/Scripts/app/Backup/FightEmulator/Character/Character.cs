using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 属性，装备，可战斗主体
    public partial class Character : IOwnerImpl, ICharacter
    {
        private string _name;
        public string name { get { return _name; } }
        public void SetName(string s)
        {
            _name = s;
        }

        private int _id = -1;
        public int id { get { return _id; } }
        public void SetId(int id) { _id = id; }

        // 场景id
        private int _sceneId = -1;
        public int sceneId { get { return _sceneId; } }
        public void SetSceneId(int id) { _sceneId = id; }

        // 位置
        public Vector3 pos = Vector3.zero;


        private int _level = 1;
        public int level { get { return _level; } }
        public void SetLevel(int level) { _level = level; }

        // 职业
        private int _class = (int)eClass.None;
        public int classType { get { return _class; } }
        public void SetClass(int classType) { _class = classType; }

        // 阵营
        private int _side = 0;
        public void SetSide(int side) { _side = side; }
        public int side { get { return _side; } }

        public float guardRange = 6.0f;
        public float attackRange = 2.0f;
       

        private Attrs _attrs = AttrsFactory.Alloc();
        public Attrs attrs { get { return _attrs; } }

        private CharAttrs _charAttrs = new CharAttrs();
        public CharAttrs charAttrs { get{ return _charAttrs; } }

        private EquipComponent _equips = new EquipComponent();

        private Skill.SkillCtrl _skill = new Skill.SkillCtrl();
        public Skill.SkillCtrl skillCtrl { get { return _skill; } }

        private Skill.PhysicAttackCtrl _PACtrl = new Skill.PhysicAttackCtrl();
        private Skill.SkillCDCtrl _cd = new Skill.SkillCDCtrl();
        private Skill.SkillSlots _skillSlots = new Skill.SkillSlots();
        public Skill.SkillSlots skillSlots { get { return _skillSlots; } }
        private Skill.SkillCastCtrl _castCtrl = new Skill.SkillCastCtrl();

        private AICtrl _ai = new AICtrl();
        private MoveCtrl _moveCtrl = new MoveCtrl();
        public MoveCtrl moveCtrl { get { return _moveCtrl; } }

        protected BagSystem.PlayerBags _bags;
        public BagSystem.PlayerBags bags { get { return _bags; } }
        public void SetBags(BagSystem.PlayerBags bags)
        {
            _bags = bags;
        }

        public Character()
        {
            _charAttrs.InitAttrs(this, _attrs);
            _equips.Init(this);
            _skill.Init(this);
            _castCtrl.Init(this, _skillSlots);
            _ai.Init(this);
            _moveCtrl.Init(this);
        }

        public int GetLevel()
        {
            return _level;
        }

        public void InitMonster(string monsterId, int level)
        {
            MonsterData ed = MonsterDataMgr.It.GetItem(monsterId);
            if (ed == null)
                return;
            _charAttrs.InitEnemy(monsterId);
            SetLevel(level);

            _charAttrs.AddFeature(EnemyFeature.It);
            PrepareAttrs();

            var weapon = GetMainWeapon();
            _PACtrl.SetHandAttackInterval(Skill.eHandType.MainHand, weapon.GetSpeed());
        }

        // 设置了基本种族，职业之后
        public void PrepareAttrs()
        {   
            // 根据职业设置feature
            prepareClassFeature();
            // 设置数值转化
            prepareClassAttrTransform();

            _charAttrs.FullyCalcAttrs();

            // hp
            _attrs.GetAttr(AttrDefine.HP).Base.baseValue 
                = _attrs.GetAttr(AttrDefine.HPMax).intFinal;
        }

        private void prepareClassFeature()
        {
            _charAttrs.AddFeature(ClassFeatureFactory.It.GetClassFeature(_class));
        }

        private void prepareClassAttrTransform()
        {
            ClassTransformerBuilder.It.BuildTransformer(_charAttrs, _class);
        }        

        public void ChangeLevel(int level)
        {
            if (level == _level)
                return;           

            var off = level - _level;
            _level = level;
            _charAttrs.Levelup(off);
        }

        public bool IsMonster()
        {
            return !string.IsNullOrEmpty(_charAttrs.GetMonsterId());
        }

        // 装备可装备
        public IWeapon GetMainWeapon()
        {
            if (IsMonster())
                return _charAttrs.GetEnemyWeapon();
            return _equips.GetSlotItem(eEquipSlot.MainHand);
        }

        public IWeapon GetOffHandWeapon()
        {
            return _equips.GetSlotItem(eEquipSlot.OffHand);
        }

        public bool CheckEquip(eEquipSlot slotIndex, EquipItem item, CheckEquipResult result)
        {   
            _equips.CheckEquipItem(slotIndex, item, result);
            return result.allOK;
        }

        public bool Equip(eEquipSlot slotIndex, int indexInBag, EquipItem item)
        {
            if (item == null)
                return false;
            return _equips.EquipItem(slotIndex, indexInBag, item);
        }

        public bool Unequip(eEquipSlot slotIndex)
        {
            return _equips.UnequipItem(slotIndex);
        }

        public void ApplyEquip(eEquipSlot slotIndex, EquipItem item)
        {
            item.Equip(_charAttrs);
            // TODO: 特性列表


            onSlotEquiped(slotIndex);
        }

        public void ApplyUnequip(eEquipSlot slotIndex, EquipItem item)
        {
            item.Unequip(_charAttrs);

            onSlotUnequiped(slotIndex);
        }

        private void onSlotEquiped(eEquipSlot slotIndex)
        {
            // 如果是武器，更新武器速度
            if(slotIndex == eEquipSlot.MainHand)
            {
                var weapon = GetMainWeapon();
                _PACtrl.SetHandAttackInterval(Skill.eHandType.MainHand, weapon.GetSpeed());
                return;
            }
            if(slotIndex == eEquipSlot.OffHand)
            {
                var weapon = GetOffHandWeapon();
                if (GetShieldBlock() == 0)
                    _PACtrl.SetHandAttackInterval(Skill.eHandType.OffHand, weapon.GetSpeed());
                else
                    _PACtrl.SetHandNoWeapon(Skill.eHandType.OffHand);
                return;
            }
        }

        private void onSlotUnequiped(eEquipSlot slotIndex)
        {
            if (slotIndex == eEquipSlot.MainHand)
            {   
                _PACtrl.SetHandNoWeapon(Skill.eHandType.MainHand);
                return;
            }
            if (slotIndex == eEquipSlot.OffHand)
            {
                _PACtrl.SetHandNoWeapon(Skill.eHandType.OffHand);
                return;
            }
        }

        public EquipItem GetSlotItem(eEquipSlot slotIndex)
        {
            return _equips.GetSlotItem(slotIndex);
        }

        public void Dump()
        {
            Log.LogCenter.Default.Debug("name: {0} lv: {1}", _name, _level);
            _attrs.Dump();
            DumpEquips();
        }

        public void DumpEquips()
        {
            _equips.Dump();
            var weapon = GetMainWeapon();
            if(weapon != null)
            {
                Log.LogCenter.Default.Debug("weapon: ({0}-{1}) speed:{2}",
                    weapon.GetMinDmg(), weapon.GetMaxDmg(), weapon.GetSpeed());
            }
        }

        public void OnAttrChanged(string attrName, float oldV, float newV)
        {
            // 可以做一些处理            
            if(attrName == AttrDefine.HPMax)
            {
                // 血量变化，保持百分比
                keepHPPercent(oldV, newV);
                if (newV < oldV)
                {
                    if (_charAttrs.GetHP() > _charAttrs.GetHPMax())
                        _charAttrs.SetHP(_charAttrs.GetHPMax());
                }
            }

            // TODO: 攻速变化,更新普攻cd
        }

        // 来回折腾由于截断原因，会有损
        private void keepHPPercent(float oldV, float newV)
        {
            float hp = _charAttrs.GetHP();
            float oldP = hp / oldV;
            float newHP = oldP * newV;

            // 避免搞成了0点血
            if (hp > 1 && newHP < 1)
                newHP = 1;

            _charAttrs.SetHP((int)newHP);
        }

        public void Update()
        {   
            autoAttack();
            _skill.Update();
            _ai.Update();
            _moveCtrl.Update();
        }

        public EquipSlot[] GetSlots()
        {
            return _equips.slots;
        }
    }
}// namespace Phoenix
