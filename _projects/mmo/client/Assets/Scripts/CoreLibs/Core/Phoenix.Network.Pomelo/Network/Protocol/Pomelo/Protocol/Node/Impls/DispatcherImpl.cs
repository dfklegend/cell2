using System;

namespace Phoenix.Network.Protocol.Pomelo
{
    // 抽象两种处理器
    // 回调
    // 接口Service
    public interface IDispatcherImpl
    {
        void SetSerializer(Serializer.ISerializer serializer);
        bool HasHandler(string route);
        void InvokeRequest(IClientSession session, 
            string route, byte[] rawData, Action<object> cbFinish);
        void InvokeNotify(IClientSession session,
            string route, byte[] rawData);
    }
}

