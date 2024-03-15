using System;
using System.Net.Sockets;

namespace Phoenix.Network
{
    public static class SocketUtil
    {
        public static void SafeClose(Socket socket)
        {
            try
            {
                socket.Shutdown(SocketShutdown.Both);
            }
            catch (Exception e)
            {
                Phoenix.Utils.SystemUtil.LogHandledException(e);
            }
            finally
            {
                socket.Close();
            }
        }
    }
}
