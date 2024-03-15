using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace PomeloCommon
{
    public class Say
    {
        public string str;
    }

    public class SayRet
    {
        public string str;
    }

    public class PushTestMsg
    {
        public string str;
        public int v;
    }    

    public class TestMsg
    {
        public string str;
    }

    public class EmptyMsg
    {
    }

    public class QueryGate 
    { 
    }

    public class QueryGateAck
    {
        public int code;
        public string ip = "";
        public string port = "";
    }

    public class LoginReq
    {
        public string Name;
        // 主动要求进入某个聊天服务器
        public int chatId = -1;
    }

    public class NormalAck
    {
        public int Code;
        public string Result;
    }

    public class ChatMsg
    {
        public string Name;
        public string Content;
    }

    // gate->chat
    public class RoomEntry
    {
        public string UId;
        public string Name;
        public string ServerId;
        public int netId;
    }

    // gate->chat
    public class RoomLeave
    {   
        public string ServerId;
        public int netId;
    }

    // chat->client
    public class OnMembers
    {
        public List<string> Members = new List<string>();
    }
}
