using Phoenix.Core;
using Phoenix.Utils;
using System.Collections.Generic;
using UnityEngine;
using System.Linq;

namespace Phoenix.Game
{	
	public class UIMgr : Singleton<UIMgr>
    {
        private Transform _root;
        private Canvas _canvas;
        private Camera _camera;
        public Camera camera { get { return _camera; } }

        private Dictionary<string, BasePanel> _panels = new Dictionary<string, BasePanel>();
        private List<string> _panelNames = new List<string>();

        StringToClassFactory<BasePanel> _factory = new StringToClassFactory<BasePanel>();
        
        private List<BasePanel> _panelsInScene = new List<BasePanel>();
        private bool _zDirt = false;

        private List<string> _toDel = new List<string>();



        public void Init()
        {
            var go = GameObject.Find("UI/MainCanvas");
            if (!go)
                return;
            _root = go.transform;
            _canvas = go.GetComponent<Canvas>();
            _camera = _canvas.worldCamera;
            _factory.RegisterAll();

            AppTask.CommonTaskDriver.It.AddTask(new UIUpdater());
        }

        public void Update()
        {
            // TODO: 关闭长时间不可见界面
            updatePanels();
        }

        public void LateUpdate()
        {            
            updateZOrder();
        }

        // 先统一路径
        private string GetAssetPath(string name)
        {
            return $"Panels/{name}";
        }

        public BasePanel OpenPanel(string name)
        {
            BasePanel panel;
            if (_panels.TryGetValue(name, out panel))
                return panel;
            panel = _factory.Create(name);
            if (panel == null)
            {
                Log.LogCenter.Default.Error("can not find panel:{0}", name);
                return null;
            }
            panelsAdd(name, panel);
            panel.OnCreate();
            panel.SetLoading();
            string assetPath = GetAssetPath(name);
            // 创建对应的panel对象
            var req = Resources.LoadAsync(assetPath, typeof(GameObject));
            req.completed += (op) => 
            {
                var prefab = req.asset as GameObject;
                if(prefab == null)
                {
                    Log.LogCenter.Default.Error("can not find prefab:{0}", assetPath);
                    return;
                }
                panel.OnLoaded(instancePrefab(req.asset as GameObject));
            };
            return panel;
        }

        public BasePanel GetPanel(string name)
        {
            BasePanel panel;
            if (_panels.TryGetValue(name, out panel))
                return panel;
            return null;
        }

        public T GetPanel<T>(string name)
            where T: BasePanel
        {
            BasePanel panel;
            if (_panels.TryGetValue(name, out panel))
                return panel as T;
            return null;
        }

        public T GetPanel<T>()
            where T : BasePanel
        {
            BasePanel panel;
            var type = typeof(T);
            var name = type.Name;
            if (_panels.TryGetValue(name, out panel))
                return panel as T;
            return null;
        }

        private void panelsAdd(string name, BasePanel panel)
        {
            _panels[name] = panel;
            _panelNames.Add(name);
        }

        private void panelsRemove(string name)
        {
            _panels.Remove(name);
            _panelNames.Remove(name);
        }

        private GameObject instancePrefab(GameObject prefab)
        {
            if (!prefab)
                return null;
            return GameObject.Instantiate(prefab);
        }

        private void addGoToScene(GameObject go)
        {
            if (!go)
                return;                 
            go.transform.SetParent(_root, false);
            var rect = go.GetComponent<RectTransform>();            

            rect.anchoredPosition = Vector2.zero;            
            go.transform.localScale = Vector3.one;            
        }

        private void removeGoFromScene(GameObject go)
        {
            if (!go)
                return;
            go.transform.SetParent(null);
        }

        public void AddToScene(BasePanel panel)
        {
            _panelsInScene.Add(panel);
            addGoToScene(panel.GetGameObject());
            SetZDirt();
        }

        public void RemoveFromScene(BasePanel panel)
        {
            _panelsInScene.Remove(panel);
            removeGoFromScene(panel.GetGameObject());
            SetZDirt();
        }

        // 面板排序
        public void SetZDirt()
        {
            _zDirt = true;
        }

        private void updateZOrder()
        {
            if (!_zDirt)
                return;
            _panelsInScene.Sort((a, b) => { return a.depth - b.depth; });
            for(int i = 0; i < _panelsInScene.Count; i ++)
            {
                var panel = _panelsInScene[i];
                if (!panel.IsReady())
                    continue;
                panel.GetGameObject().transform.SetSiblingIndex(i);
            }
        }

        private void updatePanels()
        {
            // 面板Destroy流程，保证不会在循环内移除
            // 在循环内添加会产生，但是逻辑上不会有问题(add加在后面)
            List<string> keys = _panelNames;
            for (var i = 0; i < keys.Count; i ++)
            {
                var name = keys[i];
                BasePanel panel;
                if (!_panels.TryGetValue(name, out panel))
                    continue;
                if (panel.Destroyed)
                    _toDel.Add(name);
                else
                    panel.Update();
            }

            if(_toDel.Count > 0)
            {
                for(var i = 0; i < _toDel.Count; i ++)
                {
                    panelsRemove(_toDel[i]);
                }
                _toDel.Clear();
            }
        }
    }
} // namespace Phoenix
