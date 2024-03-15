using UnityEngine;
using Phoenix.Game.FightEmulator;


namespace Phoenix.Game.Skill
{	    
    // 技能可以配置子技能
    // 满足多段击的需求
    public class Skill : ISkill
    {
        private string _skillId;
        public string skillId { get { return _skillId; } }
        private SkillData _cfgData;

        private Character _owner;
        public Character owner { get { return _owner; } }
        // 目标
        private int _tarId;
        public int tarId { get { return _tarId; } }
        // TODO: 地点目标        

        eSkillPhase _phase = eSkillPhase.Init;
        private float _timeStat;
        // 攻击速度可能影响普攻时间
        private float _totalTime;
        private float _hitTime;

        // break
        private bool _failed = false;
        public bool failed { get { return _failed; } }

        // 防止配置错误，嵌套过多
        private int _stackDepth = 0;
        public void SetStackDepth(int depth) { _stackDepth = depth; }
        private Skill _parent;
        public void SetParent(Skill skill) { _parent = skill; }
        private int _subIndexRunning = -1;
        private Skill _subRunning;
      
        private eHandType _whichHand = eHandType.None;

        private bool _onlyBalance = false;

        public bool Init(Character owner, string skillId)
        {
            _owner = owner;
            _skillId = skillId;
            _cfgData = SkillDataMgr.It.GetItem(skillId);
            if (_cfgData == null)
                return false;
            return true;
        }

        public void SetTarId(int id)
        {
            _tarId = id;
        }

        public void SetHandType(eHandType hand)
        {            
            _whichHand = hand;
        }

        public eHandType GetHandType()
        {
            return _whichHand;
        }

        public void SetOnlyBalance(bool v)
        {
            _onlyBalance = v;
        }

        public IWeapon GetCurWeapon()
        {
            if (_whichHand == eHandType.OffHand)
                return _owner.GetOffHandWeapon();
            return _owner.GetMainWeapon();
        }

        public SkillData cfgData { get { return _cfgData; } }

        public void Start()
        {
            if(_onlyBalance)
            {
                doBalance();
                return;
            }
            Log.LogCenter.Default.Debug("{0} skill {1} start", _owner.id, skillId);
            // 准备阶段
            _totalTime = _cfgData.totalTime;
            _hitTime = _cfgData.hitTime;
            
            applyTimeScale();

            Core.HEventUtil.Dispatch(Core.GlobalEvents.It.events,
                new HEventStartSkill(_owner, _skillId, _tarId));

            onStart();
            _timeStat = TimeUtil.Now();

            if (hasPrehitPhase())
                beginPrehit();
            else
                doHit();
        }

        private void doBalance()
        {
            ApplyOneTarget(_tarId);
            _phase = eSkillPhase.Over;
        }

        public void Update()
        {
            updatePhase();
        }

        private void updatePhase()
        {
            switch(_phase)
            {
                case eSkillPhase.Prehit:
                    updatePrehit();
                    break;
                case eSkillPhase.Posthit:
                    updatePosthit();
                    break;
                case eSkillPhase.SubskillRunning:
                    updateSubSkills();
                    break;
            }
        }

        private void applyTimeScale()
        {
            // 普攻考虑攻速加成
        }

        public float GetPrefireTime()
        {
            return _totalTime - _hitTime;
        }

        private void onStart()
        {
            _owner.OnSkillStart(this);
        }

        private bool hasPrehitPhase()
        {
            return _hitTime > 0;
        }

        private void beginPrehit()
        {
            _phase = eSkillPhase.Prehit;
        }

        private void updatePrehit()
        {
            if (checkBreak())
                return;
            // 检查时间
            var now = TimeUtil.Now();
            if(now >= _timeStat + _hitTime)
            {
                doHit();
            }
        }

        private bool isConditionsFail()
        {
            // 检查技能取消
            // 消耗不够
            // 目标失效
            return false;
        }

        private void doHit()
        {
            // 考虑：
            // 如果是类似wow英勇打击等附带在普攻上的技能，可以更换技能id
            Log.LogCenter.Default.Debug("{0} skill {1} dohit", _owner.id, skillId);
            applyCost();            
            applyHit();

            // 如果有subskills，会启动subskills
            beginPhaseNextHit();
        }

        private void applyCost()
        {

        }        

        private void applyHit()
        {
            // do hit to every target
            _owner.OnSkillHit(this);

            SkillImplFactory.It.ApplyHit(_cfgData.type, this);

            //_owner.TestAttackEnemy(this, _tarId);
            // 收集目标
            // 结算每一个目标
            //ApplyOneTarget(_tarId);
        }

        public void ApplyOneTarget(int tarId)
        {
            Character tar = FightCtrl.It.GetChar(tarId);
            if (tar == null)
                return;            

            FormulaResult result = FormulaUtil.Formula(this, _owner, tar) as FormulaResult;
            if (result == null || result.data.Dmg == 0)
                return;
            tar.ApplyDmg((int)result.data.Dmg);

            Log.LogCenter.Default.Debug($"{_owner.name}攻击{tar.name}，skill: {this.skillId} 造成{(int)result.data.Dmg}点伤害");
            Log.LogCenter.Default.Debug($"{tar.name}剩余血量: {tar.charAttrs.GetHP()}");

            // 抛出事件
            Core.HEventUtil.Dispatch(Core.GlobalEvents.It.events,
                new HEventAttack(_owner, tar, result));
        }

        private void beginPhaseNextHit()
        {
            if(hasSubSkills())
            {
                beginSubSkills();
                return;
            }
            beginPosthit();
        }

        private bool hasSubSkills()
        {
            return _cfgData.subSkills != null &&
                _cfgData.subSkills.Length > 0;
        }

        private void beginSubSkills()
        {
            if(_stackDepth >= 5)
            {
                Log.LogCenter.Default.Error("skill嵌套层次太多:{0}", _skillId);
                onBreak();
                return;
            }

            _phase = eSkillPhase.SubskillRunning;
            startSubSkill(0);            
        }

        private Skill createNextSubSkill(int index)
        {
            if (index >= _cfgData.subSkills.Length)
                return null;
            var subSkillId = _cfgData.subSkills[index];
            var subSkill = SubSkillCreator.CreateSubSkill(
                _owner, this, subSkillId, _stackDepth+1);
            return subSkill;
        }

        private bool hasSubSkill(int index)
        {
            return index < _cfgData.subSkills.Length;
        }

        private void updateSubSkills()
        {
            if(_subRunning.failed)
            {
                onBreak();
                return;
            }

            if (!_subRunning.IsOver())
                _subRunning.Update();

            if (_subRunning.IsOver())
            {
                if (hasSubSkill(_subIndexRunning + 1))
                    startSubSkill(_subIndexRunning + 1);
                else
                {
                    onEnd();
                }
            }   
        }

        private void startSubSkill(int index)
        {
            var subSkill = createNextSubSkill(index);
            if(subSkill == null)
            {
                onBreak();
                return;
            }

            subSkill.Start();
            _subIndexRunning = index;
            _subRunning = subSkill;
        }
       

        private void beginPosthit()
        {
            _phase = eSkillPhase.Posthit;
        }

        private void updatePosthit()
        {
            var now = TimeUtil.Now();
            if (now >= _timeStat + _totalTime)
            {
                onEnd();
            }
        }

        private bool checkBreak()
        {
            if(isConditionsFail())
            {
                onBreak();
                return true;
            }
            return false;
        }

        private void onEnd()
        {
            Log.LogCenter.Default.Debug("{0} skill {1} onEnd", _owner.id, skillId);
            SetOver();
        }

        // 技能被打断
        // 被动或者条件变化
        private void onBreak()
        {            
            Log.LogCenter.Default.Debug("{0} skill {1} onBreak", _owner.id, skillId);
            _failed = true;            
            SetOver();
            _owner.OnSkillStop(this, true);
        }

        // 主动取消
        private void onCancel()
        {
            _failed = true;
            SetOver();
            _owner.OnSkillStop(this, false);
        }

        public void SetOver()
        {
            _phase = eSkillPhase.Over;
        }

        public bool IsOver()
        {
            return _phase == eSkillPhase.Over;
        }
    }
} // namespace Phoenix
