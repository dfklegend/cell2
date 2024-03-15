using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using Phoenix.Utils;

namespace Phoenix.Log
{
    // 便于重载成其他的logger
    public interface ILogger
    {
        string GetName();
        void Debug(string format, params object[] args);
        void Info(string format, params object[] args);
        void Warning(string format, params object[] args);
        void Error(string format, params object[] args);
    }    
}