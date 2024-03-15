using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using System.Net.Sockets;


namespace Phoenix.Game.Card
{
    public class HEventLoginSucc: HEvent<int>
    {
        public long UId;
        public HEventLoginSucc(long uid)
            : base(EventDefine.LoginSucc)
        {
            UId = uid;
        }
    }

    public class HEventConnectFailed : HEvent<int>
    {
        public SocketError error;
        public HEventConnectFailed(SocketError err)
            : base(EventDefine.ConnectFailed)
        {
            error = err;
        }
    }

    public class HEventConnectBroken : HEvent<int>
    {
        public SocketError error;
        public HEventConnectBroken(SocketError err)
            : base(EventDefine.ConnectionBroken)
        {
            error = err;
        }
    }
}// namespace Phoenix
