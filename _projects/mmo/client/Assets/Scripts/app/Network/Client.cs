using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Text;
using Phoenix.API;
using Phoenix.Network;
using Phoenix.Game;
using PomeloCommon;
using Phoenix.Log;
using Phoenix.Core;
using Phoenix.Game.Card;
using System.Net.Sockets;
using Phoenix.Network.Protocol.Pomelo;

namespace Network
{
    public enum eClientState
    {
        Init,
        QueryGate,
        Connect2,
        Connect2_OK,
        Logined,
        Error,
    }

    public class Client
    {
        private PomeloClient _client;
        private int _tryTimes = 0;
        private eClientState _state = eClientState.Init;
        private bool _stopped = false;

        public string username = "";

        private PomeloClient CreateClient()
        {
            var client = new PomeloClient("client", new Phoenix.Serializer.ProtoSerializer());

            client.cbConnectFail = (error) =>
            {
                LogCenter.Default.Debug("cbConnectFail");
                // 连接失败
                onConnectFail(error);
            };

            client.cbConnectionErr = (error) =>
            {
                LogCenter.Default.Debug("cbConnectionErr");
                // 中断了
                onConnectionErr(error);
            };

            client.cbConnected = (c) =>
            {
                LogCenter.Default.Debug("cbConnected");
                onConnected();
            };
            return client;
        }

        // 从头开始
        public void Start(string ip, int port)
        {
            _stopped = false;
            if (_client != null)
            {
                _client.Close();
                _client = null;
            }

            var ipEndPoint = new IPEndPoint(IPAddress.Parse(ip), port);
            _client = CreateClient();
            _client.Connect(ipEndPoint);
            _state = eClientState.QueryGate;

            LogCenter.Default.Debug($"connect to: {ip} : {port}");
        }

        public bool IsConnected()
        {
            return _client.IsConnected();
        }

        public PomeloClient GetClient()
        {
            return _client;
        }

        public void SetState(eClientState state)
        {
            _state = state;
        }

        public void Stop()
        {
            _stopped = true;
            if (_client != null)
            {
                _client.Close();
                _client = null;
            }
            _state = eClientState.Init;
        }

        private void onConnected()
        {
            _tryTimes = 0;

            if (_state == eClientState.QueryGate)
            {
                queryGate(_client);
                return;
            }

            if (_state == eClientState.Connect2)
            {
                // begin login
                login(_client);
                _state = eClientState.Connect2_OK;
                return;
            }
        }

        public bool IsReady()
        {
            return _state == eClientState.Logined;
        }

        private void onConnectFail(SocketError error)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventConnectFailed(error));
        }

        private void onConnectionErr(SocketError error)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventConnectBroken(error));
        }

        private bool isReachMaxTryTimes()
        {
            return _tryTimes >= 5;
        }

        private void tryReConnect()
        {
            if (isReachMaxTryTimes())
            {
                LogCenter.Default.Debug("ChatClient connect faild! reach max try times");
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

            var result = await client.node.RequestAsync<Cproto.QueryGateAck>("gate.gate.querygate", new Cproto.QueryGateReq());
            if (!result.IsSucc())
            {
                LogCenter.Default.Debug($"Fatal error querygate code: {result.code}");
                return;
            }

            var ack = result.GetData<Cproto.QueryGateAck>();
            LogCenter.Default.Debug($"{ack}");
            var ports = ack.Port.Split(',');
            if (ports.Length != 2)
                return;
            _state = eClientState.Connect2;
            client.Connect(ack.IP, int.Parse(ports[1]));
        }

        private async void login(PomeloClient client)
        {
            if (!client.IsConnected())
                return;

            var result = await client.node.RequestAsync<Cproto.LoginAck>("gate.gate.login",
                new Cproto.LoginReq() { Username = username });
            if (!result.IsSucc())
            {
                LogCenter.Default.Debug($"login code: {result.code}");

                Stop();
                return;
            }

            var ack = result.GetData<Cproto.LoginAck>();
            if (ack == null)
            {
                LogCenter.Default.Debug($"bad login ack: {ack}");
                Stop();
                return;
            }

            if (ack.Code != 0)
            {
                LogCenter.Default.Debug($"bad login ack, code: {ack.Code} uid:{ack.UId}");
                Stop();

                UIMgr.It.GetPanel<PanelDialog>().ShowInfo($"登录失败，错误码: {ack.Code}");
                return;
            }


            _state = eClientState.Logined;

            LogCenter.Default.Debug($"登录成功，uid: {ack.UId}");
            UIMgr.It.GetPanel<PanelDialog>().ShowInfo($"登录成功，uid: {ack.UId}");

            // dispatch
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventLoginSucc(ack.UId));
        }

        public async void ReqCharInfo(Action cb)
        {
            var client = _client;
            if (!client.IsConnected())
                return;

            var result = await client.node.RequestAsync<Cproto.NormalAck>("scene.scene.reqcharinfo",
                new Cproto.EmptyArg() {});
            if (!result.IsSucc())
            {
                LogCenter.Default.Debug($"reqcharinfo code: {result.code}");               
                return;
            }

            var ack = result.GetData<Cproto.NormalAck>();
            if (ack == null)
            {
                LogCenter.Default.Debug($"bad reqcharinfo ack: {ack}");                
                return;
            }

            if (ack.Code != 0)
            {
                LogCenter.Default.Debug($"bad reqcharinfo ack, code: {ack.Code}");                
            }

            LogCenter.Default.Debug($"req char info succ");
            cb();
        }

        public async void OpenCamera()
        {
            var client = _client;

            if (!client.IsConnected())
                return;

            var result = await client.node.RequestAsync<Cproto.NormalAck>("scene.scene.opencamera",
                new Cproto.EmptyArg() { });
            if (!result.IsSucc())
            {
                LogCenter.Default.Debug($"opencamera code: {result.code}");
                return;
            }

            var ack = result.GetData<Cproto.NormalAck>();
            if (ack == null)
            {
                LogCenter.Default.Debug($"bad opencamera ack: {ack}");
                return;
            }

            if (ack.Code != 0)
            {
                LogCenter.Default.Debug($"bad opencamera ack, code: {ack.Code}");
            }

            LogCenter.Default.Debug($"opencamera succ");
            // 后续将触发loadscene
        }


        public async void StartFight(int downId, int upId)
        {
            var client = _client;

            if (!client.IsConnected())
                return;

            var result = await client.node.RequestAsync<Cproto.NormalAck>("logic.logic.startfight",
                new Cproto.StartGame() { DownCard = downId, UpCard = upId });
            if (!result.IsSucc())
            {
                LogCenter.Default.Debug($"startfight code: {result.code}");
                return;
            }

            var ack = result.GetData<Cproto.NormalAck>();
            if (ack == null)
            {
                return;
            }

            if (ack.Code != 0)
            {
                UIMgr.It.GetPanel<PanelDialog>().ShowInfo($"请求战斗失败，战斗还在进行中");
                LogCenter.Default.Debug($"startfight failed,  {ack.Err}");
                return;
            }


            LogCenter.Default.Debug("startfight succ");
            //ClientApp.It.stateCtrl.ChangeState((int)Phoenix.Game.eAppState.FightFE);
        }

        public async void OnLoadSceneOver()
        {
            var client = _client;

            if (!client.IsConnected())
                return;
            await client.node.RequestAsync<Cproto.NormalAck>("scene.scene.clientloadsceneover",
                new Cproto.ClientLoadSceneOver()
                {
                    UId = DataCenter.It.UId,
                    SceneId = DataCenter.It.SceneId,
                });
        }

        public async void SystemCmd<CmdResult>(string system, string cmd, object args, Action<CmdResult, int> cb)
            where CmdResult: class
        {
            var client = _client;

            if (!client.IsConnected())
                return;

            var serializer = Phoenix.Serializer.ProtoSerializer.sharedSerializer;

            var reqCmd = new Cproto.ReqSystemCmd()
                {
                    System = system,
                    Cmd = cmd,
                    
                };

            if( args != null) 
            {
                reqCmd.Args = Google.Protobuf.ByteString.AttachBytes(serializer.Serialize(args));
            }

            var result = await client.node.RequestAsync<Cproto.AckSystemCmd>("scene.scene.systemcmd",
                reqCmd);

            if (!result.IsSucc())
            {
                LogCenter.Default.Debug($"SystemCmd code: {result.code}");
                cb?.Invoke(null, (int)result.code);
                return;
            }

            var ack = result.GetData<Cproto.AckSystemCmd>();
            if (ack == null)
            {
                return;
            }

            if(ack.Code != 0)
            {
                cb?.Invoke(null, ack.Code);
                return;
            }

            if (ack.Ret == null)
            {
                cb?.Invoke(null, ack.Code);
                return;
            }
            
            var retObj = serializer.Deserialize(ack.Ret.ToByteArray(), typeof(CmdResult));
            cb?.Invoke(retObj as CmdResult, ack.Code);
        }

        public async void SystemCmd(string system, string cmd, object args, Action<int> cb)            
        {
            var client = _client;

            if (!client.IsConnected())
                return;

            var serializer = Phoenix.Serializer.ProtoSerializer.sharedSerializer;

            var reqCmd = new Cproto.ReqSystemCmd()
                {
                    System = system,
                    Cmd = cmd,
                    
                };

            if( args != null) 
            {
                reqCmd.Args = Google.Protobuf.ByteString.AttachBytes(serializer.Serialize(args));
            }

            var result = await client.node.RequestAsync<Cproto.AckSystemCmd>("scene.scene.systemcmd",
                reqCmd);

            if (!result.IsSucc())
            {
                LogCenter.Default.Debug($"SystemCmd code: {result.code}");
                cb?.Invoke((int)result.code);
                return;
            }

            var ack = result.GetData<Cproto.AckSystemCmd>();
            if (ack == null)
            {
                cb?.Invoke((int)-1);
                return;
            }

            cb?.Invoke(ack.Code);
        }        
    }    
}
