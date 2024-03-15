using System;
using System.Collections.Generic;

namespace Phoenix.Network.Protocol.Pomelo
{
    using NotifyAction = Action<IClientSession, object>;
    using RequestAction = Action<IClientSession, object, Action<object>>;

    public class RequestHandler
    {
        public Type argType;
        public RequestAction cb;
    }

    public class NotifyHandler
    {
        public Type argType;
        public NotifyAction cb;
    }

    public class FuncDispatcherImpl : IDispatcherImpl
    {
        // 只能有一个处理者
        Dictionary<string, RequestHandler> _handlers =
            new Dictionary<string, RequestHandler>();

        // 通知处理者
        Dictionary<string, NotifyHandler> _notify =
            new Dictionary<string, NotifyHandler>();

        Serializer.ISerializer _serializer;

        public void SetSerializer(Serializer.ISerializer serializer)
        {
            _serializer = serializer;
        }

        public bool HasHandler(string route)
        {
            return _handlers.ContainsKey(route) || _notify.ContainsKey(route);
        }

        private RequestHandler doRegisterRequestHandler(string route, RequestAction action)
        {
            var handler = new RequestHandler();
            handler.cb = action;
            _handlers[route] = handler;
            return handler;
        }

        public void RegisterRequestHandler(string route, RequestAction action)
        {
            doRegisterRequestHandler(route, action);
        }

        public void RegisterRequestHandler<T>(string route, RequestAction action)
        {
            var handler = doRegisterRequestHandler(route, action);
            handler.argType = typeof(T);            
        }

        private NotifyHandler doRegisterNotifyHandler(string route, NotifyAction action)
        {
            var handler = new NotifyHandler();            
            handler.cb = action;
            _notify[route] = handler;
            return handler;
        }

        public void RegisterNotifyHandler(string route, NotifyAction action)
        {
            doRegisterNotifyHandler(route, action);
        }

        public void RegisterNotifyHandler<T>(string route, NotifyAction action)
        {
            var handler = doRegisterNotifyHandler(route, action);
            handler.argType = typeof(T);            
        }

        public void InvokeRequest(IClientSession session, 
            string route, byte[] rawData, Action<object> cbFinish)
        {
            RequestHandler handler;
            _handlers.TryGetValue(route, out handler);
            if (handler == null)
            {                
                return;
            }

            // 序列化成目标结构
            object data = rawData;
            if (_serializer != null && handler.argType != null)
            {
                data = _serializer.Deserialize(rawData, handler.argType);
            }

            // 调用handler
            try
            {
                handler.cb.Invoke(session, data, cbFinish);
            }
            catch(Exception e)
            {
                Console.WriteLine("Exception:");
                Console.WriteLine(e);
            }
            
        }

        public void InvokeNotify(IClientSession session,
            string route, byte[] rawData)
        {
            NotifyHandler handler;
            _notify.TryGetValue(route, out handler);
            if (handler == null)
            {   
                return;
            }

            object data = rawData;
            if (_serializer != null && handler.argType != null)
            {
                data = _serializer.Deserialize(rawData, handler.argType);
            }            

            try
            {
                handler.cb.Invoke(session, data);
            }
            catch (Exception e)
            {
                Console.WriteLine("Exception:");
                Console.WriteLine(e);
            }
        }
    }    
}

