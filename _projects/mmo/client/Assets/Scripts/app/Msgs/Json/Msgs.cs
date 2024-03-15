using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Network
{
    public class EmptyArg
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
        public string username;
        public string password;
    }

    public class NormalAck
    {
        public int code;
        public string err;
    }

    public class LoginAck
    {
        public int code;
        public Int64 uid;
    }

    public class StartGame
    {
    }

    public class CharInfo
    {
        public string name;
        public int level;
        public Int64 exp;
        public Int64 money;
    }

    public class BattleLog
    {
        public string log;
    }
}
