using System;
using Phoenix.API;

namespace Phoenix.Network.Protocol.Pomelo
{
    // 转化为对Service.Func的调用
    public class ServiceDispatcherImpl : IDispatcherImpl
    {
        private APICollection _apis;

        public ServiceDispatcherImpl(string categories)
        {
            _apis = CollectionBuilder.BuildFromService(categories, null);
        }

        public void SetSerializer(Serializer.ISerializer serializer)
        {
            _apis.SetSerializer(serializer);
        }        

        bool IDispatcherImpl.HasHandler(string route)
        {
            return _apis.HasAPI(route);
        }

        private IContext toContext(IClientSession session)
        {
            var p = session as PomeloSession;
            return p;
        }

        public void InvokeRequest(IClientSession session,
            string route, byte[] rawData, Action<object> cbFinish)
        {           
            _apis.InvokeRequest(route, toContext(session), rawData, cbFinish);
        }

        public void InvokeNotify(IClientSession session,
            string route, byte[] rawData)
        {
            _apis.InvokeNotify(route, toContext(session), rawData);
        }
    }
}

