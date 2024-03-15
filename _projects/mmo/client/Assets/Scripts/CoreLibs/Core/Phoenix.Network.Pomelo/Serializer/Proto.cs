using System;
using System.Text;
using Google.Protobuf;

namespace Phoenix.Serializer
{
    public class ProtoSerializer : ISerializer
    {
        static byte[] _empty = Encoding.UTF8.GetBytes("");
        public static ISerializer sharedSerializer = new ProtoSerializer();
        public byte[] Serialize(object value)
        {
            try 
            {
                var message = (IMessage)value;                
                return message.ToByteArray();                
            }
            catch(Exception e)
            {
                Console.WriteLine("Exception:");
                Console.WriteLine(e);
                Phoenix.Utils.SystemUtil.LogHandledException(e);
                return _empty;
            }            
        }

        public object Deserialize(byte[] value, Type type)
        {
            try
            {
                Object message = Activator.CreateInstance(type);
                ((IMessage)message).MergeFrom(value);
                return message;
            }
            catch(Exception e)
            {
                Console.WriteLine("Exception:");
                Console.WriteLine(e);
                Phoenix.Utils.SystemUtil.LogHandledException(e);
                return null;
            }
        }
    }
}

