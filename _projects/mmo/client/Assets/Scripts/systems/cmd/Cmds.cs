using System;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Utils;
using System.Text;

[StringType("test")]
public class GMCmdTest: BaseGMCmd
{
    
    public override void Execute(string[] args, Action<bool> cb)
    {
        GMCmdUtils.Output($"test cmd executed");        
    }
}

// 发送命令到logic执行
[StringType("lcmd")]
public class GMCmdLogicCmd : BaseGMCmd
{

    public async override void Execute(string[] args, Action<bool> cb)
    {
        var client = ClientApp.It.client;
        if(!client.IsConnected())
        {
            GMCmdUtils.Output("未连接");
            return;
        }

        // 组装命令
        StringBuilder sb = new StringBuilder();
        for (var i = 1; i < args.Length; i++)
        {
            if (i > 1)
                sb.Append(" ");
            sb.Append(args[i]);
        }

        var result = await client.GetClient().node.RequestAsync<Cproto.CmdAck>("scene.scene.cmd",
            new Cproto.Cmd()
            {
                Cmd_ = sb.ToString()
            });
        if (!result.IsSucc())
        {
            cb?.Invoke(false);
            return;
        }

        var ack = result.GetData<Cproto.CmdAck>();
        if(ack == null)
        {
            return;
        }

        GMCmdUtils.Output(ack.Result);
        cb?.Invoke(true);
    }
}