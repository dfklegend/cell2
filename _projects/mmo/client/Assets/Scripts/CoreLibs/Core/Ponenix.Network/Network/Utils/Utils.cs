using System.Net;
using System.Threading.Tasks;

namespace Phoenix.Network
{
    public static class NetworkUtils
    {
        public static IPAddress SafeParse(string address)
        {
            IPAddress ret;
            if (IPAddress.TryParse(address, out ret))
                return ret;
            return IPAddress.Parse("127.0.0.1");
        }

        // 
        public static IPEndPoint tryParsePoint(string address)
        {
            string[] subs = address.Split(':');
            if (subs.Length < 2)
                return null;
            int port = 0;
            int.TryParse(subs[1], out port);
            return new IPEndPoint(SafeParse(subs[0]), port);
        }

        // 避免直接调用异步函数的警告
        public static void NoWarnning(this Task task)
        {
        }
    }
}

