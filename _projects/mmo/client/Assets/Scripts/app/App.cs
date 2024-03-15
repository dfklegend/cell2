using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Client;
using Phoenix.Log;
using Phoenix.Core;
using Benchtest;
using Phoenix.Scheduler;
using System.Threading;
using System.Threading.Tasks;
using Phoenix.Network;
using Phoenix.Utils;
using Phoenix.Network.Protocol.Pomelo;
using Phoenix.Game;
using Network;

public class ClientApp : BaseClientApp 
{
    private static ClientApp _it;
    public static ClientApp It
    {
        get
        {
            return _it;
        }
    }   


    AppStateCtrl _fsm; 
    public AppStateCtrl stateCtrl 
    {
        get 
        {
            return _fsm;
        }
    }

    Client _client = new Client();
    public Client client
    {
        get { return _client; }
    }

    public ClientApp()
    {
        _it = this;
    }

    public override void OnPrepare()
    {
        base.OnPrepare();

        ThreadMgrConfig config = new ThreadMgrConfig();
        config.skipNewMainContext = true;
        ThreadMgr.Start(2, config);
        AppComponentMgr.Register<AppModule>();
        Phoenix.Game.Card.DataCenter.It.Start();

        // rpc性能测试
        //new BenchTest().Start();
    }

    public override void OnReady()
    {
        base.OnReady();

        LogCenter.GetLogInfo("Trace").SetConsolePrintingEnabled(false);
        LogCenter.Default.Debug("OnReady");

        _fsm = FSMFactory.Create();
        _fsm.ChangeState(_fsm.factory.GetInitState());
    }    

    public override void OnUpdate()
    {
        base.OnUpdate();
        ThreadMgr.it.mainThread.PullUpdate();
        _fsm.Update();
    }
    
    public void OnQuit()
    {
        _client.Stop();
    }
}
