using System.Collections.Generic;
using Newtonsoft.Json.Linq;

namespace Phoenix.Core.Json
{
    // 根节点
    public class NodeEntry: BaseEntry
    {
        protected Dictionary<string, BaseEntry> _entries
            = new Dictionary<string, BaseEntry>();

        public void AddEntry(string name, BaseEntry entry)
        {
            _entries[name] = entry;
        }

        public T GetEntry<T>(string name)
            where T: BaseEntry
        {
            BaseEntry ret;
            if (_entries.TryGetValue(name, out ret))
                return ret as T;
            return null;
        }

        public override void LoadFromJson(JToken node)
        {
            if (node == null)
                return;
            foreach (var key in _entries.Keys)
            {
                var entry = _entries[key];
                JToken jObj = node[key];
                if (jObj == null)
                    continue;
                entry.LoadFromJson(jObj);
            }
        }

        public override void SaveToJson(JToken node)
        {
            if (node == null)
                return;
            foreach (var key in _entries.Keys)
            {
                var entry = _entries[key];

                var jObj = new JObject();
                node[key] = jObj;

                entry.SaveToJson(jObj);
            }
        }

        public void FromStr(string str)
        {
            JObject root = JsonUtil.SafeParse(str);
            LoadFromJson(root);
        }

        public string ToStr()
        {
            JObject root = new JObject();
            SaveToJson(root);
            return root.ToString();
        }
    }
}
