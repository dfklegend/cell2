using System;
using System.Collections.Generic;

namespace Phoenix.API
{
    // 一个接口集合
    public class APICollection
    {
        private Dictionary<string, APIEntry> _entries = new Dictionary<string, APIEntry>();
        // 用来序列化参数
        Serializer.ISerializer _serializer;

        public APICollection(Serializer.ISerializer argSerializer)
        {
            _serializer = argSerializer;
        }

        public void SetSerializer(Serializer.ISerializer serializer)
        {
            _serializer = serializer;
            foreach(var one in _entries.Values)
            {
                one.SetSerializer(serializer);
            }
        }

        public APIEntry GetEntry(string name, bool createIfMiss = true)
        {
            APIEntry entry;
            if (_entries.TryGetValue(name, out entry))
                return entry;
            if (!createIfMiss)
                return null;
            entry = new APIEntry(_serializer);
            _entries[name] = entry;
            return entry;
        }

        // 如果route只有一个
        public static string[] SplitRoute(string route)
        {
            string[] subs = route.Split('.');
            if (subs.Length >= 2)
                return subs;
            // 只有一个
            // 转换为对 系统接口的调用
            if (subs.Length == 1)
            {
                string[] newSubs = new string[2];
                newSubs[0] = "";
                newSubs[1] = subs[0];
                return newSubs;
            }
            // 0
            return subs;
        }

        public bool HasAPI(string route)
        {
            string[] subs = SplitRoute(route);
            if (subs.Length != 2)
                return false;
            var entry = GetEntry(subs[0], false);
            if (entry == null)
                return false;
            return entry.HasAPI(subs[1]);
        }

        private (string, string, APIEntry) prepareEntry(string route)
        {
            string[] subs = SplitRoute(route);
            if (subs.Length != 2)
            {
                return (string.Empty, string.Empty, null);
            }

            return (subs[0], subs[1], GetEntry(subs[0], false));
        }

        // service(entry name).apiname
        public bool InvokeRequest(string route, IContext context, byte[] arg, Action<object> cbFinish)
        {
            var (entryName, methodName, entry) = prepareEntry(route);
            if (entry == null)
            {
                cbFinish(null);
                return false;
            }
            if(!entry.InvokeRequest(methodName, context, arg, cbFinish))
            {
                cbFinish(null);
                return false;
            }
            return true;
        }

        public bool InvokeNotify(string route, IContext context, byte[] arg)
        {
            var (entryName, methodName, entry) = prepareEntry(route);
            if (entry == null)
                return false;
            return entry.InvokeNotify(methodName, context, arg);
        }
    }
}

