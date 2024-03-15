using System;
using System.Text;

namespace Phoenix.Network.Protocol.Pomelo
{
    // 定义一些公共变量
    public static class PomeloProtocolGlobal
    {
        // 缺省定义的dispatcher
        public static RichMsgDispatcher msgDispatcher = new RichMsgDispatcher();

        public static void Init()
        {
            msgDispatcher.RegisterNotifyHandler("__error__", (session, data) => {
                LogHandlerError(data);
            });
        }

        public static void LogHandlerError(object data)
        {
            byte[] rawData = (byte[])data;
            var err = Encoding.UTF8.GetString(rawData);
            Console.WriteLine("error: " + err);
        }

        public static void LogHandlerError(string data)
        {   
            Console.WriteLine("error: " + data);
        }

        // 处理消息超时
        public static void Update()
        {
            msgDispatcher.Update();
        }
    }   
}

