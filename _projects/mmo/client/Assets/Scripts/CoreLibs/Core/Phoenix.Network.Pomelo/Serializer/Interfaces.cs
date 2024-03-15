using System;

namespace Phoenix.Serializer
{

    // 序列化库
    // 根据一个目标类型来序列化 
    public interface ISerializer
    {
        byte[] Serialize(object value);
        object Deserialize(byte[] value, Type type);
    }
}

