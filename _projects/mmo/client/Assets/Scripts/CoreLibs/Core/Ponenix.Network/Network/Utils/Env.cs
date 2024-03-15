using Phoenix.Log;

namespace Phoenix.Network
{
    // 模块环境，提供log支持
    public static class Env
    {
        static LogProxy _log = new LogProxy();
        // 可以定制Log
        public static LogProxy logProxy { get { return _log; } }
        public static ILogger L {get { return _log.L; } }
        
    }    
}

