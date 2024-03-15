#define SIMPLEJSON

using System;
using System.Text;

namespace Phoenix.Serializer
{
#if SIMPLEJSON
    public class SimpleJsonSerializer : ISerializer
    {
        static byte[] _emptyJson = Encoding.UTF8.GetBytes("{}");
        public byte[] Serialize(object value)
        {
            try 
            {
                var str = SimpleJson.SimpleJson.SerializeObject(value);
                return Encoding.UTF8.GetBytes(str);
            }
            catch(Exception e)
            {
                Console.WriteLine("Exception:");
                Console.WriteLine(e);
                Phoenix.Utils.SystemUtil.LogHandledException(e);
                return _emptyJson;
            }            
        }

        public object Deserialize(byte[] value, Type type)
        {
            try
            {
                return SimpleJson.SimpleJson.DeserializeObject(Encoding.UTF8.GetString(value), type);
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
#endif
}

