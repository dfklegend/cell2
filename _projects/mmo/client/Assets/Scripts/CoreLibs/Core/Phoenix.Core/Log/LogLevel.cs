using System;
using System.ComponentModel;

namespace Phoenix.Log
{
    public enum LogLevel
    {
        [Description("None")]
        None = 0,
        [Description("Debug")]
        Debug = 1,        
        [Description("Info")]
        Info = 2,        
        [Description("Warning")]
        Warning = 3,
        [Description("Error")]
        Error = 4
    }
}

