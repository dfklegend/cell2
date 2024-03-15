using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using Phoenix.Utils;

namespace Phoenix.Core
{
    // ID分配器
    // TODO: 并发服务器如果为了保证唯一性，可能需要有个
    // 服务器节点
    // > 0
    public interface IUniqueIDAllocer
    {
        // 分配ID
        ulong AllocId();
    }

    public class SimpleUniqueIDAllocer : IUniqueIDAllocer
    {
        private ulong _nextId = 1;
        public ulong AllocId()
        {
            return _nextId ++;
        }
    }

    // 有一个全局分配，分配前缀
    // 本地有部分可以自增长分配
    // [前缀][自增长]
    public class ServerUniqueIDAllocer : IUniqueIDAllocer
    {
        public ulong AllocId()
        {
            return 0;
        }
    }
}