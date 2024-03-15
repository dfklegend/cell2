using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using System.Collections.Generic;

namespace Phoenix.Game.Card
{	
    public class UnitModelMgr : Singleton<UnitModelMgr>
    {
        private Transform _root;
        private Dictionary<int, BaseModel> _models = new Dictionary<int, BaseModel>();
        private List<BaseModel> _modelsForSort = new List<BaseModel>();
        
        private GameObject _prefabChar;
        private GameObject _prefabBullet;
        private GameObject _prefabExit;

        public void Init(Transform root)
        {
            _root = root;
            BindEvents(true);

            _prefabChar = Resources.Load<GameObject>("Panels/prefabs/char");
            _prefabBullet = Resources.Load<GameObject>("Panels/prefabs/bullet");
            _prefabExit = Resources.Load<GameObject>("Panels/prefabs/exit");
        }

        public void Clear()
        {
            BindEvents(false);
            // destroy all
            foreach( var kv in _models)
            {
                kv.Value.Destroy();
            }
            _models.Clear();
            _modelsForSort.Clear();
        }

        private void BindEvents(bool bind)
        {
            var events = GlobalEvents.It.events;
            events.Bind(EventDefine.InitUnit, OnInitUnit, bind);
            events.Bind(EventDefine.HPChanged, OnHPChanged, bind);
            events.Bind(EventDefine.MPChanged, OnMPChanged, bind);
            events.Bind(EventDefine.StartSkill, OnStartSkill, bind);
            events.Bind(EventDefine.SkillBroken, onSkillBroken, bind);
        }

        private void OnInitUnit(params object[] args)
        {
            var e = args[0] as HEventInitUnit;           
            UnitModel model = getUnitModel(e.unitId);
            if (model == null)
                return;
            var c = FightCtrl.It.GetChar(e.unitId);
            model.InitInfo(c);
        }

        private UnitModel getUnitModel(int id)
        {
            BaseModel bm;
            if (!_models.TryGetValue(id, out bm))
                return null;
            return bm as UnitModel;
        }

        private void OnStartSkill(params object[] args)
        {
            var e = args[0] as HEventStartSkill;
            UnitModel model = getUnitModel(e.src.id);
            if (model == null)
                return;
            Log.LogCenter.Default.Debug("{0} onStartSkill:{1}", 
                e.src.id, e.skillId);
            var sd = SkillDataMgr.It.GetItem(e.skillId);
            if (sd != null && !string.IsNullOrEmpty(sd.action))
            {             
                //model.PlayAnim(sd.action, 0.5f);
                model.PlayAnim(sd.action, sd.totalTime/1000f);
                model.LookAtTar(e.tarId);
            }
        }

        private void onSkillBroken(params object[] args)
        {
            var e = args[0] as HEventSkillBroken;
            UnitModel model = getUnitModel(e.src.id);
            if (model == null)
                return;
            Log.LogCenter.Default.Debug("{0} onSkillBroken:{1}",
                e.src.id, e.skillId);
            var sd = SkillDataMgr.It.GetItem(e.skillId);
            if (sd != null && !string.IsNullOrEmpty(sd.action))
            {
                model.StopCurAnim();
            }
        }

        private void OnHPChanged(params object[] args)
        {
            var e = args[0] as HEventHPChanged;
            UnitModel model = getUnitModel(e.src.id);
            if (model == null)
                return;
            model.UpdateHP(e.src.GetHPPercent());
        }

        private void OnMPChanged(params object[] args)
        {
            var e = args[0] as HEventMPChanged;
            UnitModel model = getUnitModel(e.src.id);
            if (model == null)
                return;
            model.UpdateMP(e.src.GetMPPercent());
        }

        private GameObject createGo()
        {
            return GameObject.Instantiate(_prefabChar);
        }

        private GameObject createBulletGo()
        {
            return GameObject.Instantiate(_prefabBullet);
        }

        private GameObject createExitGo()
        {
            return GameObject.Instantiate(_prefabExit);
        }

        public UnitModel CreateModel(CharacterUnit unit)
        {
            UnitModel model = new UnitModel();
            var go = createGo();
            go.transform.SetParent(_root, false);

            model.Init(go, unit);

            _models[unit.entity.GetEntityID()] = model;
            _modelsForSort.Add(model);            
            return model;
        }

        public BulletModel CreateBullet(BulletUnit unit)
        {
            BulletModel model = new BulletModel();
            var go = createBulletGo();
            go.transform.SetParent(_root, false);

            model.Init(go, unit);

            _models[unit.entity.GetEntityID()] = model;
            _modelsForSort.Add(model);            
            return model;
        }

        public StaticModel CreateExit(StaticUnit unit)
        {
            var model = new StaticModel();
            var go = createExitGo();
            go.transform.SetParent(_root, false);

            model.Init(go, unit);

            _models[unit.entity.GetEntityID()] = model;
            _modelsForSort.Add(model);

            return model;            
        }

        public UnitModel GetModel(int id)
        {
            return getUnitModel(id);
        }

        public BaseModel GetBaseModel(int id)
        {
            BaseModel bm;
            if (!_models.TryGetValue(id, out bm))
                return null;
            return bm;
        }

        public void UpdateModelPos(int id)
        {
            var model = GetModel(id);
            if (model == null)
                return;
            model.UpdatePos();
        }

        public void Update()
        {
            var values = _modelsForSort;
            foreach(var one in values)
            {
                one.Update();
            }
            UpdateZOrder();
        }

        private void UpdateZOrder()
        {            
            _modelsForSort.Sort((a, b) =>
            {                
                return (int)((a.depth - b.depth) * 10);
            });

            for (int i = 0; i < _modelsForSort.Count; i++)
            {
                var model = _modelsForSort[i];
                model.GetTransform().SetSiblingIndex(i);
            }
        }

        public void DestroyModel(int id)
        {
            var model = GetModel(id);
            if (model == null)
                return;
            model.Destroy();

            _models.Remove(id);
            _modelsForSort.Remove(model);
        }

        public void DestroyBaseModel(int id)
        {
            var model = GetBaseModel(id);
            if (model == null)
                return;
            model.Destroy();

            _models.Remove(id);
            _modelsForSort.Remove(model);
        }
    }    
} // namespace Phoenix
