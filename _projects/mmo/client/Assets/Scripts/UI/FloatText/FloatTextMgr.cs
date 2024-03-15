using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game.FText
{    
    public class FloatTextMgr : Singleton<FloatTextMgr>
    {
        Transform _root;
        private List<FloatText> _texts = new List<FloatText>();
        GameObject _prefab;

        public void Init(GameObject go)
        {
            if (!go)
                return;
            _root = go.transform;

            _prefab = Resources.Load<GameObject>("Panels/prefabs/floattext1");
        }     
        
        public void Update()
        {
            for(var i = 0; i < _texts.Count;)
            {
                var one = _texts[i];
                if (one.IsOver())
                {
                    one.Destroy();
                    _texts.RemoveAt(i);
                }
                else
                {
                    one.Update();
                    i++;
                }
            }
        }

        private GameObject createGo() => GameObject.Instantiate(_prefab);

        public FloatText Create()
        {
            FloatText text = new FloatText();
            var go = createGo();
            
            go.transform.SetParent(_root, false);
            text.Init(go);

            _texts.Add(text);
            return text;
        }
    }
} // namespace Phoenix
