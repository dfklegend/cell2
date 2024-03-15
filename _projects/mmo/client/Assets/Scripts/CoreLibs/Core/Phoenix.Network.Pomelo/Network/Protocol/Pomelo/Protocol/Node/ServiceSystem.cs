using Phoenix.API;

namespace Phoenix.Network.Protocol.Pomelo
{
    [APIService("__system__", "")]
    public class ServiceSystem: IAPIService
    {
        [APIFunc("__error__")]
        public void ErrorHandle(IContext context, string data)
        {
            PomeloProtocolGlobal.LogHandlerError(data);
        }
    }
}
