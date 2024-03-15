using System;
using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;

namespace Phoenix.Game
{ 
    public class SystemCmdItem
    {
        public Action<object> cb;
        public Type argsType;
    }

    // 系统可以接收服务器的本系统的Cmd，用来更新数据
    // 系统可以发起Request，用来向服务器对应system发起请求
    public abstract class BaseSystem
    {
        private Dictionary<string, SystemCmdItem> _cmds = new Dictionary<string, SystemCmdItem>();

        public abstract string GetName();
        public virtual void Reset() { }
        
        public void OnCmd(string cmd, byte[] argsBuf)
        {
            SystemCmdItem item;
            if(!_cmds.TryGetValue(cmd, out item))            
            {
                Log.LogCenter.Default.Debug($"{GetName()} has not cmd: {cmd}");
                return;
            }
                
            var serializer = Serializer.ProtoSerializer.sharedSerializer;

            if(item.argsType == null)
            {
                item.cb.Invoke(null);
                return;
            }
            var data = serializer.Deserialize(argsBuf, item.argsType);
            item.cb.Invoke(data);
        }

        protected void registerCmd<T>(string cmd, Action<object> cb)
        {
            var item = new SystemCmdItem();
            item.argsType = typeof(T);
            item.cb = cb;
            _cmds[cmd] = item;
        }

        public void Request<CmdResult>(string cmd, object args, Action<CmdResult, int> cb)
            where CmdResult : class
        {
            var client = ClientApp.It.client;
            client.SystemCmd(GetName(), cmd, args, cb);
        }

        public void Request(string cmd, object args, Action<int> cb)          
        {
            var client = ClientApp.It.client;
            client.SystemCmd(GetName(), cmd, args, cb);
        }
    }    
}
