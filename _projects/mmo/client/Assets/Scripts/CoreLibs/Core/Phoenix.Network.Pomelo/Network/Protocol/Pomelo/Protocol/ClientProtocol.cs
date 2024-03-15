using System;
using System.Text;
using Pomelo.DotNetClient;
using SimpleJson;
using TimeUtil = Phoenix.Utils.TimeUtil;

namespace Phoenix.Network.Protocol.Pomelo
{
    public enum eClientState
    {
        Init,
        WaitHandshake,
        Working,
        Closed
    }

    // 根据服务器的心跳间隔时间，发心跳包
    public class ClientProtocol : BasePomeloProtocol
    {
        protected eClientState _state = eClientState.Init;

        private float _nextHeartBeat = 0f;       

        protected void setState(eClientState state)
        {
            _state = state;
        }

        protected bool isState(eClientState state)
        {
            return _state == state;
        }

        public override void OnConnected()
        {
            base.OnConnected();
            setState(eClientState.WaitHandshake);
            sendClientHandShake();            
        }

        public override bool IsReady()
        {
            return isState(eClientState.Working);
        }

        public override void Update()
        {
            if (IsStopped())
                return;
            base.Update();            
            trySendHeartbeat();
        }

        private void trySendHeartbeat()
        {
            if (_heartBeatInterval > 0)
            {
                var now = TimeUtil.Now();
                if (now >= _nextHeartBeat)
                {
                    _nextHeartBeat = now + _heartBeatInterval;
                    send(PackageType.PKG_HEARTBEAT);
                }
            }
        }

        private void sendClientHandShake()
        {
            //Console.WriteLine($"sendClientHandShake");
            byte[] body = Encoding.UTF8.GetBytes(buildMsg(null).ToString());
            send(PackageType.PKG_HANDSHAKE, body);
        }

        private JsonObject buildMsg(JsonObject user)
        {
            if (user == null) user = new JsonObject();

            JsonObject msg = new JsonObject();

            //Build sys option
            JsonObject sys = new JsonObject();
            sys["version"] = PomeloDefine.Version;
            sys["type"] = PomeloDefine.Type;

            //Build handshake message
            msg["sys"] = sys;
            msg["user"] = user;
            return msg;
        }

        protected override void processHandshake(PomeloMsg msg)
        {
            if (!isState(eClientState.WaitHandshake))
            {

                Env.L.Error("got handshake msg when state is not waithandshake");
                return;
            }

            //Console.WriteLine($"client.processHandshake");

            JsonObject data = (JsonObject)SimpleJson.SimpleJson.DeserializeObject(
                Encoding.UTF8.GetString(msg.package.body));
            processHandshakeData(data);
        }

        private void processHandshakeData(JsonObject msg)
        {
            //Handshake error
            if (!msg.ContainsKey("code") || !msg.ContainsKey("sys") || Convert.ToInt32(msg["code"]) != 200)
            {
                //throw new Exception("Handshake error! Please check your handshake config.");
                Env.L.Error("Handshake error! Please check your handshake config.");
                return;
            }

            //Set compress data
            JsonObject sys = (JsonObject)msg["sys"];

            JsonObject dict = new JsonObject();
            if (sys.ContainsKey("dict")) dict = (JsonObject)sys["dict"];

            JsonObject protos = new JsonObject();
            JsonObject serverProtos = new JsonObject();
            JsonObject clientProtos = new JsonObject();

            if (sys.ContainsKey("protos"))
            {
                protos = (JsonObject)sys["protos"];
                serverProtos = (JsonObject)protos["server"];
                clientProtos = (JsonObject)protos["client"];
            }

            _messageProtocol = new MessageProtocol(dict, serverProtos, clientProtos);

            //Init heartbeat service
            int interval = 0;
            if (sys.ContainsKey("heartbeat")) interval = Convert.ToInt32(sys["heartbeat"]);
            
            _heartBeatInterval = interval;       
            if(interval > 0)
            {
                _nextHeartBeat = TimeUtil.Now() + interval;
            }

            //send ack and change protocol state            
            send(PackageType.PKG_HANDSHAKE_ACK, new byte[0]);

            setState(eClientState.Working);
            OnReady();
            //Console.WriteLine($"client.OnReady");

            //Invoke handshake callback
            JsonObject user = new JsonObject();
            if (msg.ContainsKey("user")) user = (JsonObject)msg["user"];
            //handshake.invokeCallback(user);
        } 

        protected override void processKick(PomeloMsg msg)
        {
        }
    }

    public class ClientProtocolFactory : IProtocolFactory
    {
        public IProtocol Create()
        {
            return new ClientProtocol();
        }
    }
}

