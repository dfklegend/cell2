using System;
using SimpleJson;

namespace Pomelo.DotNetClient
{
    public class Message
    {
        public MessageType type;
        public string route;
        public uint id;
        public bool error;

        public JsonObject oldJsonData;
        // 原始数据
        public byte[] rawData;
        // 最后反序列化
        public object data;

        public Message(MessageType type, uint id, string route, JsonObject data)
        {
            this.type = type;
            this.id = id;
            this.route = route;
            this.oldJsonData = data;
            this.error = false;
        }

        public Message(MessageType type, uint id, string route, byte[] data, bool err)
        {
            this.type = type;
            this.id = id;
            this.route = route;
            this.rawData = data;
            this.error = err;
        }
    }
}