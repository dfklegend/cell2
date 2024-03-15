using System;
using Phoenix.Core;
using Phoenix.Utils;

public class GMCmdMgr : Singleton<GMCmdMgr>
{
    StringToClassFactory<BaseGMCmd> _factory = new StringToClassFactory<BaseGMCmd>();

    GMCmdMgr()
    {
        _factory.RegisterAll();
    }

    private BaseGMCmd CreateCmd(string op)
    {
        return _factory.Create(op);
    }

    public void Execute(string input, Action<bool> cb)
    {
        string[] subs = input.Split(' ');

        var cmd = CreateCmd(subs[0]);
        if (cmd == null) 
        {
            GMCmdUtils.Output($"找不到命令: {input}");
            return;
        }         
        cmd.Execute(subs, cb);
    }
}