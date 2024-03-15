using System;
using System.IO;

namespace Phoenix.Network
{
    public interface IMsg
    {        
    }

    public interface IMsgEncoder
    {
        // 写入stream
        Stream Write(IMsg msg);
    }

    public interface IMsgDecoder
    {
        // 解析成Msg
        IMsg Make(Stream stream);
    }    

    public interface IMsgCoderFactory
    {
        IMsgEncoder CreateEncoder();
        IMsgDecoder CreateDecoder();
    }

    public interface ISession
    {
        void SendImmediately();
        void Stop();
        void ForceStop();
        int GetHandle();
    }

    public interface IChannel
    {
        int GetHandle();
        // 获得数据
        // 不能失败
        void OnGotData(byte[] data, int offset, int count);

        /*
         * Return
         *      实际返回的数据量
         */
        int PullSendData(byte[] buf, int offset, int max);
        int GetSendDataSize();
        void SendMsg(Stream s);
        void SendMsg(byte[] bytes);

        void Update();
        void Clear();

        void StopSession();
    }   

    // 可以注入一些定制数据
    public interface IConnectionImpl
    {      
    }

    public interface IProtocolConfig
    {
    }
    
    // 协议
    public interface IProtocol
    {
        int GetHandle();
        void Init(IChannel channel,
            IMsgEncoder encoder, IMsgDecoder decoder,
            IMsgProcessor processor, IProtocolConfig config);
        void OnConnected();
        // 是否已经可以处理用户协议
        bool IsReady();
        void SetCBOnReady(Action cb);
        void MakeMsg(Stream stream);
        void SendMsg(IMsg msg);
        void Update();

        void Stop();

        void SetImpl(IConnectionImpl impl);
        T GetImpl<T>() where T: class, IConnectionImpl;
    }

    public interface IProtocolFactory
    {
        IProtocol Create();
    }

    // 处理
    public interface IMsgProcessor
    {
        void Process(IProtocol protocol, IMsg msg);
        void Update();
        IMsg Pop();
    }

    public interface IProcessorFactory
    {
        IMsgProcessor Create();
    }

    public interface IConnStateFactor
    {
        IConnectionImpl CreateServerState();
        IConnectionImpl CreateClientState();
    }
}

