using System;
using System.Collections.Generic;
using System.Reflection;

namespace Phoenix.API
{
    public enum eAPIType
    {
        Request = 0,
        Notify
    }

    public class MethodInfo
    {
        public eAPIType apiType;
        public object obj;
        // 参数类型
        public Type contextType;
        public Type argType;
        public System.Reflection.MethodInfo method;
    }

    // Request(context, object arg, Action<object> cb)
    // Notify(context, object arg)
    // . 分析某个接口符合条件的函数
    // . 根据参数分析，序列化参数
    // . 调用指定的函数
    public class APIEntry
    {
        //
        // TODO: 使用RuntimeMethodHandle等
        Dictionary<string, MethodInfo> _apis = new Dictionary<string, MethodInfo>();
        Serializer.ISerializer _serializer;

        private static Type _typeFuncAttr = typeof(APIFunc);
        private static Type _typeCBFinish = typeof(Action<object>);
        private static Type _typeBaseContextType = typeof(IContext);

        public APIEntry(Serializer.ISerializer serializer)
        {
            _serializer = serializer;
        }        

        public void SetSerializer(Serializer.ISerializer serializer)
        {
            _serializer = serializer;
        }

        public void AnalysisType(Type type)
        {
            var objThis = Activator.CreateInstance(type, null);
            var typeFuncAttr = _typeFuncAttr;
            var methods = type.GetMethods();
            
            foreach(var method in methods)
            {
                if (!method.IsDefined(typeFuncAttr))
                    continue;
                // 检查一下参数
                var (correct, apiType) = checkParameters(method);
                if (!correct)
                    continue;

                var args = method.GetParameters();

                var one = new MethodInfo();
                one.obj = objThis;                
                one.apiType = apiType;
                one.contextType = args[0].ParameterType;
                one.argType = args[1].ParameterType;
                one.method = method;

                addAPI(one);
            }
        }
        
        internal (bool, eAPIType) checkParameters(System.Reflection.MethodInfo method)
        {
            var args = method.GetParameters();
            if (3 == args.Length)
            {
                return (checkRequestParameters(args), eAPIType.Request);
            }

            return (checkNotifyParameters(args), eAPIType.Notify);
        }

        // context, arg, cbFinish
        internal bool checkRequestParameters(ParameterInfo[] args)
        {   
            if (args.Length != 3)
            {
                return false;
            }

            if (!checkContext(args[0]))
                return false;

            var cbType = args[2].ParameterType;
            if (cbType != _typeCBFinish)
            {
                return false;
            }
            return true;
        }

        // context, arg
        internal bool checkNotifyParameters(ParameterInfo[] args)
        {   
            if (args.Length != 2)
            {
                return false;
            }
            if (!checkContext(args[0]))
                return false;
            return true;
        }

        internal bool checkContext(ParameterInfo context)
        {
            var type = context.ParameterType;
            var baseContextType = _typeBaseContextType;
            if (baseContextType.IsAssignableFrom(type))
                return true;
            return false;
        }

        void addAPI(MethodInfo api)
        {
            var name = APIFunc.GetName(api.method);
            _apis[name] = api;
        }

        public bool HasAPI(string name)
        {
            return _apis.ContainsKey(name);
        }

        public bool InvokeRequest(string name, IContext context, byte[] arg, Action<object> cbFinish)
        {
            MethodInfo api;
            if (!_apis.TryGetValue(name, out api))
            {
                Console.WriteLine($"{name} is not found!");
                return false;
            }
            if(api.apiType != eAPIType.Request)
            {
                Console.WriteLine($"{name} is not request!");
                return false;
            }
            object[] args = new object[3];
            args[0] = context;
            if(_serializer != null)
                args[1] = _serializer.Deserialize(arg, api.argType);
            args[2] = cbFinish;
            invokeMethod(api, args);
            return true;
        }

        public bool InvokeNotify(string name, IContext context, byte[] arg)
        {
            MethodInfo api;
            if (!_apis.TryGetValue(name, out api))
                return false;
            if (api.apiType != eAPIType.Notify)
            {
                Console.WriteLine($"{name} is not notify!");
                return false;
            }
            object[] args = new object[2];
            args[0] = context;
            if(_serializer != null)
                args[1] = _serializer.Deserialize(arg, api.argType);
            invokeMethod(api, args);
            return true;
        }

        private void invokeMethod(MethodInfo api, object[] args)
        {
            try 
            {
                api.method.Invoke(api.obj, args);
            }
            catch(Exception e)
            {
                APIUtils.LogException(e);
            }
            
        }
    }
}

