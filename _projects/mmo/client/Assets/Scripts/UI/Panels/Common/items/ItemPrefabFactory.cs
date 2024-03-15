using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

using PrefabDict = System.Collections.Generic.Dictionary<Phoenix.Game.eItemStyle, UnityEngine.GameObject>;
namespace Phoenix.Game
{    
    public class ItemPrefabFactory : Singleton<ItemPrefabFactory>
    {

        // 列表里面的模板
        //private Dictionary<eItemStyle, GameObject> _prefabs = new Dictionary<eItemStyle, GameObject>();
        private PrefabDict _prefabs = new PrefabDict();
        // icon形态的模板
        //private Dictionary<eItemStyle, GameObject> _icons = new Dictionary<eItemStyle, GameObject>();
        private PrefabDict _icons = new PrefabDict();

        public void Init()
        {
            if (_prefabs.Count > 0)
                return;
            initListPrefabs();
            initIconPrefabs();
        }

        // 列表显示模式
        private void initListPrefabs()
        {
            registerListStyle(eItemStyle.Normal, "normal");
            registerListStyle(eItemStyle.Label, "label");
            registerListStyle(eItemStyle.Equip, "normal");
            registerListStyle(eItemStyle.Skill, "normal");            
        }

        // icon显示模式
        private void initIconPrefabs()
        {
            registerIconStyle(eItemStyle.Normal, "normal");
        }


        public void register(PrefabDict prefabs, eItemStyle item, string path)
        {
            prefabs[item] = Resources.Load<GameObject>(path);
        }

        public void registerListStyle(eItemStyle item, string name)
        {
            register(_prefabs, item, getListStylePath(name));
        }

        private string getListStylePath(string name)
        {
            return $"Panels/ItemStyles/Bag/{name}";
        }

        public void registerIconStyle(eItemStyle item, string name)
        {
            register(_icons, item, getIconStylePath(name));
        }

        private string getIconStylePath(string name)
        {
            return $"Panels/ItemStyles/Icon/{name}";
        }

        private GameObject getListPrefab(eItemStyle item)
        {
            return getPrefab(_prefabs, item);
        }

        private GameObject getPrefab(PrefabDict prefabs, eItemStyle item)
        {
            GameObject prefab;
            if (prefabs.TryGetValue(item, out prefab))
                return prefab;
            return null;
        }

        public GameObject ListInstantiate(eItemStyle item)
        {
            return instantiate(_prefabs, item);
        }

        public GameObject IconInstantiate(eItemStyle item)
        {
            return instantiate(_icons, item);
        }

        public GameObject instantiate(PrefabDict prefabs, eItemStyle item)
        {
            var prefab = getPrefab(prefabs, item);
            if (prefab == null)
                return null;
            return GameObject.Instantiate(prefab);
        }
    }
} // namespace Phoenix
