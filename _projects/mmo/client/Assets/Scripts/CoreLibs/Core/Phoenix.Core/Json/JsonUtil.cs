using System.Collections.Generic;
using Newtonsoft.Json.Linq;
using System;
using System.Reflection;

namespace Phoenix.Core.Json
{
    // 依次组织成节点树
    public static class JsonUtil
    {
        public static string ToJsonStr(BaseEntry entry)
        {
            JObject root = new JObject();
            entry.SaveToJson(root);
            return root.ToString();
        }

        public static void FromJsonStr(BaseEntry entry, string str)
        {
            try 
            {
                JObject root = JObject.Parse(str);
                entry.LoadFromJson(root);
            }
            catch(Exception e)
            {
                PConsole.Error(e);
            }
        }


        // 载入一个list
        public static void LoadList<T>(List<T> result, JToken parent, string name)
            where T: BaseEntry, new()
        {
            var node = parent[name];
            if (node == null)
                return;
            var array = node as JArray;
            if (array == null)
                return;
            for(var i = 0; i < array.Count; i ++)
            {
                var one = new T();
                one.LoadFromJson(array[i]);
                result.Add(one);
            }            
        }

        public static void SaveList<T>(List<T> entries, JToken parent, string name)
            where T : BaseEntry, new()
        {
            var array = new JArray();
            parent[name] = array;
            
            for (var i = 0; i < entries.Count; i++)
            {
                var one = entries[i];
                var obj = new JObject();
                one.SaveToJson(obj);
                array.Add(obj);
            }
        }

        // [{k: v:},...]
        public static void LoadMap<T>(Dictionary<string, T> result, JToken parent, string name)
            where T : BaseEntry, new()
        {
            var node = parent[name];
            if (node == null)
                return;
            var array = node as JArray;
            if (array == null)
                return;
            for (var i = 0; i < array.Count; i++)
            {
                var item = array[i];
                
                var one = new T();
                one.LoadFromJson(item["v"]);

                string k = AsStr(item, "k", "");
                result.Add(k, one);
            }
        }

        
        // 转化为list来存
        public static void SaveMap<T>(Dictionary<string, T> entries, JToken parent, string name)
            where T : BaseEntry, new()
        {
            var array = new JArray();
            parent[name] = array;
            
            foreach(var key in entries.Keys)
            {
                var item = new JObject();
                var obj = entries[key];

                item["k"] = key;
                item["v"] = new JObject();
                obj.SaveToJson(item["v"]);

                array.Add(item);
            }
        }

        // 利用反射，载入T中的属性
        public static void LoadValues<T>(T obj, JToken node)
        {
            FieldInfo[] fields = obj.GetType().GetFields();
            foreach(var oneField in fields)
            {
                if (!oneField.FieldType.IsValueType &&
                     oneField.FieldType != typeof(string))
                    continue;
                if (node[oneField.Name] == null)
                    continue;
                var value = node[oneField.Name];

                try 
                {
                    // 可能版本升级造成数据类型不一致
                    oneField.SetValue(obj, value.ToObject(oneField.FieldType));
                }
                catch(Exception e)
                {
                    PConsole.Error("LoadValues exception(may be entry field type changed):");
                    PConsole.Error(e);
                }
            }
        }        

        // 利用反射，存储数值对象
        public static void SaveValues<T>(T obj, JToken node)
        {
            FieldInfo[] fields = obj.GetType().GetFields();
            foreach (var oneField in fields)
            {
                if (!oneField.FieldType.IsValueType &&
                    oneField.FieldType != typeof(string))
                    continue;
                object value = oneField.GetValue(obj);
                node[oneField.Name] = JToken.FromObject(value);
            }
        }

        public static JObject SafeParse(string str)
        {   
            try
            {
                return JObject.Parse(str);
            }
            catch(Exception e)
            {
                Utils.SystemUtil.LogHandledException(e);
            }
            return new JObject();
        }

        // 数据获取
        public static string AsStr(JToken parent, string name, string def)
        {
            var node = parent[name];
            if (node == null)
                return def;
            var v = (string)node;            
            return v!=null?v:def;
        }

        public static int AsInt(JToken parent, string name, int def)
        {
            var node = parent[name];
            if (node == null)
                return def;
            try
            {
                return (int)node;                
            }
            catch(Exception e)
            {
                PConsole.Error("error get value:"+ name);
                PConsole.Error(e);
            }
            return def;
        }
    }
}
