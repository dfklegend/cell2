using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using System.Net.Sockets;
using Network;


namespace Phoenix.Game.Card
{       
    public class HEventServerMsg<T> : HEvent<int>
    {
        public T msg;
        public HEventServerMsg(T msg, int eventType)
            : base(eventType)
        {
            this.msg = msg;
        }
    }    
}// namespace Phoenix
