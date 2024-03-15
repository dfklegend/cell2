using System;
using System.Reflection;

namespace Phoenix.API
{
    // 定义APIEntry接口类
    // [APIEntryAttr("category", "name")]
    [AttributeUsage(AttributeTargets.Class, AllowMultiple = false, Inherited = false)]
    public class APIService : Attribute
    {
        public string category { get; set; }
        public string name { get; set; }

        // 所属分类，接口名字
        public APIService(string category, string name)
        {
            this.category = category;
            this.name = name;
        }

        public static string GetName(Type t)
        {
            APIService attribute = tryGetAttr(t);
            if (attribute == null)
                return string.Empty;
            return attribute.name;
        }

        public static string GetCategory(Type t)
        {           
            APIService attribute = tryGetAttr(t);
            if (attribute == null)
                return string.Empty;
            return attribute.category;
        }

        private static APIService tryGetAttr(Type t)
        {
            return AttrUtil.TryGetAttr<APIService>(t);
        }
    }

    // 标签，说明是API函数
    [AttributeUsage(AttributeTargets.Method, AllowMultiple = false, Inherited = false)]
    public class APIFunc : Attribute
    {
        private string _name;
        public APIFunc()
        {            
        }

        public APIFunc(string name)
        {
            _name = name;
        }

        public static string GetName(MemberInfo t)
        {
            APIFunc attribute = AttrUtil.TryGetAttr<APIFunc>(t);
            if (attribute == null || string.IsNullOrEmpty(attribute._name))
                return t.Name;
            return attribute._name;
        }
    }

    public static class AttrUtil
    {
        public static T TryGetAttr<T>(Type t)
            where T : Attribute
        {
            if (!t.IsDefined(typeof(T), false))
                return null;
            return Attribute.GetCustomAttribute(t, typeof(T)) as T;
        }

        public static T TryGetAttr<T>(MemberInfo t)
            where T : Attribute
        {
            if (!t.IsDefined(typeof(T), false))
                return null;
            return Attribute.GetCustomAttribute(t, typeof(T)) as T;
        }
    }
}

