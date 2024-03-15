using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using Phoenix.Utils;
using Phoenix.Core;

/*
 * UData/log.config增加配置来配置某个log项
 * 
 * 比如:
 * [Normal]
    ConsolePrinting=true
    ScreenPrinting=true
    FilePrinting=true
    LogLevel=Debug
 * 
 *  接口:
 *  输出log
 *  Log.Asset.Debug("hello world.{0}",1);
 *  
 *  
 *  How to use
 *  add new log type
 *  1. 配置文件中配置log的参数
 *  2. 创建一个logger来方便使用
 *     public static Logger Normal = new Logger("Normal");
 *     
 *  use log
 *  Log.Asset.Debug();
 *  Log.Asset.Error();
 *  
    
 *  

//*/
namespace Phoenix.Log
{
    public class LogCenter
    {
        // 性能分析的时候关闭
        public bool Enabled = true;

        private const string CONFIG_FILE_NAME = "log.config";
        
        // 缺省就有的log分类
        private readonly LogInfo[] DEFAULT_LOG_INFOS = new LogInfo[] { 
        new LogInfo("Default", true, false, true, LogLevel.Debug),
        new LogInfo("Trace", true, false, true, LogLevel.Debug),
        new LogInfo("Asset", true, false, true, LogLevel.Debug),
        new LogInfo("Exception", true, false, true, LogLevel.Debug)};

        public static Logger Default = new Logger("Default");
        public static Logger Asset = new Logger("Asset");
        public static Logger Trace = new Logger("Trace");
        // 异常
        public static Logger Exception = new Logger("Exception");
        private Dictionary<string, Logger> m_loggers = new Dictionary<string, Logger>();

        private const string OUTPUT_DIRECTORY_NAME = "Logs";
        private const string OUTPUT_FILE_EXTENSION = "log";
        private Dictionary<string, LogInfo> m_logInfos = new Dictionary<string, LogInfo>();
        private static LogCenter s_instance;

        // 进程启动时，尽早调用
        public static void Create(string logPrefix)
        {
            LogEnv.SetLogPrefix(logPrefix);
            s_instance = new LogCenter();
            s_instance.Initialize();
        }

        public static LogCenter Get()
        {
            if (s_instance == null)
            {
                Create("");
            }
            return s_instance;
        }        

        // 获取logger，不存在则创建
        public static Logger GetLogger(string name)
        {
            var it = Get();
            Logger logger;
            if (it.m_loggers.TryGetValue(name, out logger))
                return logger;
            logger = new Logger(name);
            it.m_loggers[name] = logger;
            return logger;
        }

        private void Initialize()
        {            
            this.initDefaultInfos();
            m_loggers["Default"] = Default;
        }

        public static LogInfo GetLogInfo(string name, bool createIfMiss = false)
        {
            return Get().getLogInfo(name, createIfMiss);
        }

        private LogInfo getLogInfo(string name, bool createIfMiss = false)
        {
            LogInfo info;
            if (this.m_logInfos.TryGetValue(name, out info))                
                return info;
            if (!createIfMiss)
                return null;
            info = new LogInfo(name);            
            this.m_logInfos.Add(info.GetName(), info);
            return info;
        }        

        public bool IsConsolePrintingEnabled(string name)
        {
            if (!this.Enabled)
                return false;
            LogInfo info;
            if (!this.m_logInfos.TryGetValue(name, out info))
            {
                return false;
            }
            return info.IsConsolePrintingEnabled();
        }

        public bool IsFilePrintingEnabled(string name)
        {
            if (!this.Enabled)
                return false;
            LogInfo info;
            if (!this.m_logInfos.TryGetValue(name, out info))
            {
                return false;
            }
            return info.IsFilePrintingEnabled();
        }

        public bool IsNeedPrint(string name)
        {
            return IsConsolePrintingEnabled(name) || IsFilePrintingEnabled(name);
        }

        public bool IsPrintingEnabled(string name)
        {
            if (!this.Enabled)
                return false;
            LogInfo info;
            if (!this.m_logInfos.TryGetValue(name, out info))
            {
                return false;
            }
            return info.IsPrintingEnabled();
        }

        public bool IsScreenPrintingEnabled(string name)
        {
            if (!this.Enabled)
                return false;
            LogInfo info;
            if (!this.m_logInfos.TryGetValue(name, out info))
            {
                return false;
            }
            return info.IsScreenPrintingEnabled();
        }
        
        public void LoadConfigFile()
        {
            string path = string.Format("{0}/{1}", LogEnv.GetDataPath(), CONFIG_FILE_NAME);
            if (File.Exists(path))
            {                
                IniUtil.ParseIniFile(path, new IniUtil.ConfigFileEntryParseCallback(this.OnConfigFileEntryParsed));
            }
            
        }

        private void initDefaultInfos()
        {
            foreach (LogInfo info in this.DEFAULT_LOG_INFOS)
            {
                if (!this.m_logInfos.ContainsKey(info.GetName()))
                {
                    this.m_logInfos.Add(info.GetName(), info);
                }
            }
        }

        // 配置覆盖顺序
        // . 缺省配置
        // . 内部配置
        // . 外部配置
        public void InitConfig(string content)
        {
            if(!string.IsNullOrEmpty(content))
                IniUtil.ParseConfigText(content, new IniUtil.ConfigFileEntryParseCallback(this.OnConfigFileEntryParsed));
            LoadConfigFile();
        }

        public static bool ForceBool(string strVal)
        {
            string str = strVal.ToLowerInvariant().Trim();
            if ((!(str == "on") && !(str == "1")) && !(str == "true"))
            {
                return false;
            }
            return true;
        }

        private void OnConfigFileEntryParsed(string baseKey, string subKey, string val, object userData)
        {
            LogInfo info;
            if (!this.m_logInfos.TryGetValue(baseKey, out info))
            {
                info = new LogInfo(baseKey);
                this.m_logInfos.Add(info.GetName(), info);
            }
            if (subKey.Equals("ConsolePrinting", StringComparison.OrdinalIgnoreCase))
            {
                info.SetConsolePrintingEnabled(ForceBool(val));
            }
            else if (subKey.Equals("ScreenPrinting", StringComparison.OrdinalIgnoreCase))
            {
                info.SetScreenPrintingEnabled(ForceBool(val));
            }
            else if (subKey.Equals("FilePrinting", StringComparison.OrdinalIgnoreCase))
            {
                info.SetFilePrintingEnabled(ForceBool(val));
            }
            else if (subKey.Equals("LogLevel", StringComparison.OrdinalIgnoreCase))
            {
                try
                {
                    LogLevel logLevel = EnumUtil.GetEnum<LogLevel>(val, StringComparison.OrdinalIgnoreCase);
                    info.SetLogLevel(logLevel);
                }
                catch (ArgumentException)
                {
                }
            }
        }

        public void Print(string name, string message, LogLevel logLevel)
        {
            try
            {
                implPrint(name, message, logLevel);
            }
            catch (System.Exception e)
            {
                PConsole.Log(e);
            }
        }

        public bool CanLog(string name, LogLevel logLevel)
        {
            LogInfo info;
            return this.m_logInfos.TryGetValue(name, out info) && (logLevel >= info.GetLogLevel() );
        }

        void implPrint(string name, string message, LogLevel logLevel)
        {
            LogInfo info;
            if (this.m_logInfos.TryGetValue(name, out info) && (info.GetLogLevel() <= logLevel))
            {
                if (info.IsFilePrintingEnabled() && !info.DidFileCreateFail())
                {
                    StreamWriter fileWriter = info.GetFileWriter();
                    if (fileWriter == null)
                    {
                        string logspath = LogEnv.GetFullNameInDataPath(OUTPUT_DIRECTORY_NAME);
                        bool exists = Directory.Exists(logspath);
                        if (!exists)
                        {
                            exists = Directory.CreateDirectory(logspath).Exists;
                        }
                        if (exists)
                        {
                            try
                            {
                                fileWriter = new StreamWriter(new FileStream(string.Format("{0}/{1}.{2}",
                                    logspath, info.GetFileName(), OUTPUT_FILE_EXTENSION), FileMode.Create, FileAccess.ReadWrite));
                                info.SetFileWriter(fileWriter);
                            }
                            catch (Exception)
                            {
                                info.SetFileCreateFailed();
                            }
                        }
                        else
                        {
                            info.SetFileCreateFailed();
                        }
                    }
                    if (fileWriter != null)
                    {
                        StringBuilder builder = new StringBuilder();
                        switch (logLevel)
                        {
                            case LogLevel.Debug:
                                builder.Append("D ");
                                break;

                            case LogLevel.Info:
                                builder.Append("I ");
                                break;

                            case LogLevel.Warning:
                                builder.Append("W ");
                                break;

                            case LogLevel.Error:
                                builder.Append("E ");
                                break;
                        }
                        builder.Append(DateTime.Now.TimeOfDay.ToString());
                        builder.Append(" ");
                        builder.Append(message);
                        fileWriter.WriteLine(builder.ToString());
                        fileWriter.Flush();
                    }
                }
                if (info.IsConsolePrintingEnabled())
                {
                    string str2 = string.Format("[{0}] [{1}] {2}",
                        DateTime.Now.TimeOfDay.ToString(), name, message);
                    switch (logLevel)
                    {
                        case LogLevel.Debug:
                        case LogLevel.Info:
                            PConsole.Log(str2);
                            break;

                        case LogLevel.Warning:
                            PConsole.Warning(str2);
                            break;

                        case LogLevel.Error:
                            PConsole.Error(str2);
                            break;
                    }
                }
            }
        }
       
        public void ResetFile(string name)
        {
             LogInfo info;
             if (this.m_logInfos.TryGetValue(name, out info))
             {
                 if (info.IsFilePrintingEnabled() && !info.DidFileCreateFail())
                 {
                    StreamWriter fileWriter = info.GetFileWriter();
                    if (fileWriter != null)
                    {
                        fileWriter.Close();
                        info.SetFileWriter(null);

                        string logspath = LogEnv.GetFullNameInDataPath(OUTPUT_DIRECTORY_NAME);
                        bool exists = Directory.Exists(logspath);
                        if (!exists)
                        {
                            exists = Directory.CreateDirectory(logspath).Exists;
                        }
                        if (exists)
                        {
                            try
                            {
                                fileWriter = new StreamWriter(new FileStream(string.Format("{0}/{1}.{2}",
                                    logspath, name, OUTPUT_FILE_EXTENSION), FileMode.Create, FileAccess.ReadWrite));
                                info.SetFileWriter(fileWriter);
                            }
                            catch (Exception)
                            {
                                info.SetFileCreateFailed();
                            }
                        }
                        else
                        {
                            info.SetFileCreateFailed();
                        }
                    }
                 }
             }
        }

        public void ScreenPrint(string name, string message, LogLevel logLevel)
        {            
            // reserved
        }
    }
}

