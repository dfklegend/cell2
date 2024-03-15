using System;
using Pomelo.DotNetClient;

namespace Phoenix.Network.Protocol.Pomelo
{
    using DispatchCB = Action<IClientSession, Message>;
    public class CBMsgDispatcher : IMsgDispatcher
    {
        DispatchCB _cb;
        public CBMsgDispatcher(DispatchCB cb)
        {
            _cb = cb;
        }

        public Serializer.ISerializer GetSerializer()
        {
            return null;
        }

        public void OnSessionStart(IClientSession session)
        {
        }

        public void OnSessionStop(IClientSession session)
        {

        }

        public void Process(IClientSession session, Message msgIn)
        {
            _cb.Invoke(session, msgIn);
        }
    }    
}

