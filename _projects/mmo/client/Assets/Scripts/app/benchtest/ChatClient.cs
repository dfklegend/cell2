using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Text;
using Phoenix.API;
using Phoenix.Network;
using PomeloCommon;
using Phoenix.Log;
using Phoenix.Network.Protocol.Pomelo;

namespace Benchtest
{
    public enum eChatClientState
    {
        Init,
        QueryGate,
        Connect2,
        Connect2_OK,
        Logined
    }

    public class ChatClient
    {
        private PomeloClient _client;
        private int _tryTimes = 0;        
        private eChatClientState _state = eChatClientState.Init;
        private int _chatId = -1;

        private PomeloClient CreateClient()
        {   
            var client = new PomeloClient("client");            

            client.cbConnectFail = (error) => {
                LogCenter.Default.Debug("cbConnectFail");
                // 连接失败
                onConnectFail();
            };

            client.cbConnectionErr = (error) => {
                LogCenter.Default.Debug("cbConnectionErr");
                // 中断了
                onConnectionErr();
            };

            client.cbConnected = (c) => {
                LogCenter.Default.Debug("cbConnected");
                onConnected();
            };
            return client;
        }

        public void SetChatId(int id)
        {
            _chatId = id;
        }

        // 从头开始
        public void Start(string ip, int port)
        {
            if (_client != null)
            {
                _client.Close();
                _client = null;
            }

            var ipEndPoint = new IPEndPoint(IPAddress.Parse(ip), port);
            _client = CreateClient();
            _client.Connect(ipEndPoint);
            _state = eChatClientState.QueryGate;

            LogCenter.Default.Debug($"connect to: {ip} : {port}");
        }

        private void onConnected()
        {
            _tryTimes = 0;

            if(_state == eChatClientState.QueryGate)
            {
                queryGate(_client);
                return;
            }

            if(_state == eChatClientState.Connect2)
            {
                // begin login
                login(_client);
                _state = eChatClientState.Connect2_OK;
                return;
            }
        }

        public bool IsReadyForTest()
        {
            return _state == eChatClientState.Logined;
        }

        private void onConnectFail()
        {
            if(isReachMaxTryTimes())
            {
                // do something
            }
            tryReConnect();
        }

        private void onConnectionErr()
        {
            
            if(_state == eChatClientState.Connect2)
            {
                tryReConnect();
            }

            // 
            if (_state == eChatClientState.Logined)
            {
                _state = eChatClientState.Connect2;
                tryReConnect();
                return;
            }
        }

        private bool isReachMaxTryTimes()
        {
            return _tryTimes >= 99999;
        }

        private void tryReConnect()
        {
            if(isReachMaxTryTimes())
            {
                LogCenter.Default.Debug("ChatClient connect faild!");
                return;
            }
            _tryTimes++;

            LogCenter.Default.Debug("ChatClient connect faild!");
            LogCenter.Default.Debug("try connect again!");
            _client.ReConnect();
        }

        private async void queryGate(PomeloClient client)
        {
            if (!client.IsConnected())
                return;

            var result = await client.node.RequestAsync<QueryGateAck>("gate.gate.querygate", new EmptyMsg());
            if (!result.IsSucc())
            {
                LogCenter.Default.Debug($"Fatal error querygate code: {result.code}");
                return;
            }            

            var ack = result.GetData<QueryGateAck>();
            LogCenter.Default.Debug($"{ack}");
            var ports = ack.port.Split(',');
            if (ports.Length != 2)
                return;
            _state = eChatClientState.Connect2;
            client.Connect(ack.ip, int.Parse(ports[1]));
        }

        public async void Hello(Action<RequestResult> cb)
        {
            var client = _client;
            if (!client.IsConnected())
                return;

            var result = await client.node.RequestAsync<NormalAck>("chat.chat.hello",
                new TestMsg() { str = "hello world" });
            cb?.Invoke(result);
            if (!result.IsSucc())
                return;
            //Console.WriteLine($"chat.server.test ret: {result.GetData<NormalAck>().Result}");
        }

        private async void login(PomeloClient client)
        {
            if (!client.IsConnected())
                return;

            var result = await client.node.RequestAsync<NormalAck>("gate.gate.login",
                new LoginReq() { Name = "xx", chatId = _chatId });
            if (!result.IsSucc())
            {
                LogCenter.Default.Debug($"login code: {result.code}");

                tryReConnect();
                return;
            }

            _state = eChatClientState.Logined;
        }

        public async void Say(string str)
        {
            var client = _client;
            if (!client.IsConnected())
                return;

            var result = await client.node.RequestAsync<NormalAck>("chat.server.say",
                new ChatMsg() { Content = str });
            if (!result.IsSucc())
            {
                LogCenter.Default.Debug($"Say failed code: {result.code}");
                return;
            }
            var ack = result.GetData<NormalAck>();
            if (ack.Code != 0)
            {
                LogCenter.Default.Debug($"Say failed : {ack.Result}");
            }
        }

        public async void ReqMembers()
        {
            var client = _client;
            if (!client.IsConnected())
                return;

            var result = await client.node.RequestAsync<NormalAck>("chat.server.reqMembers",
                new EmptyMsg());
            if (!result.IsSucc())
                return;
        }
    }
    
}
