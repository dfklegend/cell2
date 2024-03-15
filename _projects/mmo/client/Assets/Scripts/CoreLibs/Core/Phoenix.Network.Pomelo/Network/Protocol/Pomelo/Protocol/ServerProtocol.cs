using System;
using System.Text;
using Pomelo.DotNetClient;
using SimpleJson;
using TimeUtil = Phoenix.Utils.TimeUtil;

namespace Phoenix.Network.Protocol.Pomelo
{
    public enum eServerState
    {
        Init,
        WaitHandshake,
        WaitHandshakeAck,
        Working,
        Closed
    }

    // 心跳超时，关闭连接
    public class ServerProtocol : BasePomeloProtocol
    {
        protected eServerState _state = eServerState.Init;

        // 过期时间，如果超时没收到timeout，断开连接
        protected float _heartBeatTimeout = 0;

        public ServerProtocol()
        {
            setState(eServerState.WaitHandshake);
            _heartBeatInterval = PomeloDefine.HEARTBEAT_INTERVAL;
        }

        protected void setState(eServerState state)
        {
            _state = state;
        }

        protected bool isState(eServerState state)
        {
            return _state == state;
        }

        public override void OnConnected()
        {
            base.OnConnected();            
        }

        public override bool IsReady()
        {
            return isState(eServerState.Working);
        }

        public override void Update()
        {
            if (IsStopped())
                return;
            base.Update();
            if(checkHeartbeatTimeout())
            {
                // 心跳超时
                ReqStopSession();
                // 避免重新调用
                updateHeartbeatTimeout();
                return;
            }            
        }

        protected bool checkHeartbeatTimeout()
        {
            if (_heartBeatTimeout == 0)
                return false;
            if (TimeUtil.Now() > _heartBeatTimeout)
            {
                Env.L.Warning("heartbeat timeout");
                // 断开连接
                return true;
            }
            return false;
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
            if (!isState(eServerState.WaitHandshake))
            {
                Env.L.Error($"server.processHandshake not in WaitHandshake");
                return;
            }

            //Console.WriteLine($"server.processHandshake");
            JsonObject data = (JsonObject)SimpleJson.SimpleJson.DeserializeObject(
                Encoding.UTF8.GetString(msg.package.body));
            processHandshakeData(data);
        }

        private void processHandshakeData(JsonObject msg)
        {
            // TODO: 可以检查下客户端版本号等

            // 发送信息给客户端
            JsonObject root = new JsonObject();
            root["code"] = 200;

            JsonObject sys = new JsonObject();
            sys["heartbeat"] = PomeloDefine.HEARTBEAT_INTERVAL;

            root["sys"] = sys;
            // dict
            // protos

            byte[] body = Encoding.UTF8.GetBytes(root.ToString());
            send(PackageType.PKG_HANDSHAKE, body);

            setState(eServerState.WaitHandshakeAck);
        }

        protected override void processHandshakeAck(PomeloMsg msg)
        {
            if (!isState(eServerState.WaitHandshakeAck))
            {
                Env.L.Error($"not in WaitHandshakeAck");
                return;
            }
            _messageProtocol = new MessageProtocol(new JsonObject(), new JsonObject(), new JsonObject());

            updateHeartbeatTimeout();
            setState(eServerState.Working);
            OnReady();            
        }

        protected override void processHeartbeat(PomeloMsg msg)
        {            
            updateHeartbeatTimeout();
            // send one back
            send(PackageType.PKG_HEARTBEAT);
        }

        protected void updateHeartbeatTimeout()
        {
            _heartBeatTimeout = TimeUtil.Now() + 3 * PomeloDefine.HEARTBEAT_INTERVAL;
        }

        protected override void processKick(PomeloMsg msg)
        {
        }
    }

    public class ServerProtocolFactory : IProtocolFactory
    {
        public IProtocol Create()
        {
            return new ServerProtocol();
        }
    }
}

