using System;

namespace Phoenix.Network.Protocol.Pomelo
{
    public class PomeloDefine
    {
        public const int LENGTH_SIZE = sizeof(UInt32);
        public const string Version = "0.3.0";
        public const string Type = "unity-socket";

        public const int HEARTBEAT_INTERVAL = 15;
        
        public static byte[] empteBytes = new byte[0];
    }   
}

