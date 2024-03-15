using System;
using System.IO;
using Phoenix.Log;

namespace Phoenix.Log
{
    // 定义给模块使用
    // 缺省指定到Default
    // 可以通过CreateCustomLogger来重新定义
    public class LogProxy
    {        
        private ILogger _logger;   
        public ILogger L { get { return _logger; } }

        public LogProxy()
        {
            _logger = LogCenter.Default;
        }       

        // 使用一个CustomLogger
        public void CreateCustomLogger(string logName, Action<LogInfo> setCB)
        {
            setCB.Invoke(LogCenter.GetLogInfo(logName, true));            
            _logger = LogCenter.GetLogger(logName);
        }

        public void CreateCustomLogger(string logName)
        {
            CreateCustomLogger(logName, (info) => {
                info.SetConsolePrintingEnabled(true)
                .SetFilePrintingEnabled(true);
            });
        }
        
        // 直接使用logger
        public void UseLogger(ILogger l)
        {
            _logger = l;
        }    
        
        public void SetLogLevel(LogLevel level)
        {
            LogCenter.GetLogInfo(_logger.GetName(), true)
                .SetLogLevel(level);
        }
    }    
}

