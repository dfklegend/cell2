using System;
using Pomelo.DotNetClient;
using SimpleJson;

namespace Phoenix.Network.Protocol.Pomelo
{
    public static class PomeloUtil
    {
        public static void DumpMsg(Message msg)
        {
            Console.WriteLine(SimpleJson.SimpleJson.SerializeObject(msg));
        }

        public static JsonObject DeserializeObject(string str)
        {
            try 
            {
                return (JsonObject)SimpleJson.SimpleJson.DeserializeObject(str);
            }
            catch(Exception e)
            {
                Console.WriteLine($"str: {str}");
                Console.WriteLine("Got exeception:");
                Phoenix.Utils.SystemUtil.LogHandledException(e);
                return new JsonObject();
            }
        }

        public static PomeloNode GetNode(IClientSession s)
        {
            PomeloSession ps = s as PomeloSession;
            if (ps == null)
                return null;
            return ps.node;
        }

        public static IClientSession GetSession(TCPConnection con)
        {
            return con.GetProtocol<BasePomeloProtocol>().node.session;
        }

        public static NetConfig InitPomeloClientCB(NetConfig config)
        {
            var pconfig = new PomeloConfig();
            
            pconfig.WithSerialize(new Serializer.SimpleJsonSerializer());            

            config.SetCoderFactory(new PomeloCoderFactory())
                .SetProtocolFactory(new ClientProtocolFactory())
                .WithProtocolConfig(pconfig)
                .SetProcessorFactory(new EmptyProcessorFactory());
            return config;
        }

        public static NetConfig InitPomeloServer(NetConfig config, IMsgDispatcher dispatcher)
        {
            var pconfig = new PomeloConfig();

            pconfig.WithMsgDispatcher(dispatcher);
            pconfig.WithSerialize(new Serializer.SimpleJsonSerializer());            

            config.SetCoderFactory(new PomeloCoderFactory())
                .SetProtocolFactory(new ServerProtocolFactory())
                .WithProtocolConfig(pconfig)
                .SetProcessorFactory(new EmptyProcessorFactory());
            return config;
        }
    }
}

