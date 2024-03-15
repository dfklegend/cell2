using System;

namespace Phoenix.Utils
{
    public static class LogUtil
    {
        public static void CreateCustomLogger(string logName, Action<Log.ILogger> postAction)
        {
            var logger = Log.LogCenter.GetLogger(logName);
            Log.LogCenter.GetLogInfo(logName, true)
                .SetConsolePrintingEnabled(true)
                .SetFilePrintingEnabled(true);
           
            postAction?.Invoke(logger);
        }
    }
} // Phoenix.Utils