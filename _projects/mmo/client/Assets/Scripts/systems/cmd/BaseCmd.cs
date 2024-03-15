using System;

public abstract class BaseGMCmd
{
    // 命令被拆解
    public abstract void Execute(string[] args, Action<bool> cb);
}

