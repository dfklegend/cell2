using System;
namespace Phoenix.Log
{
    public class Logger : ILogger
    {
        private string m_name;
        private bool _ForceLog;

        public Logger(string name)
        {
            this.m_name = name;
            _ForceLog = false;
        }

        public string Name { get { return m_name; } }

        public string GetName() { return this.m_name; }

        public void SetForceLog(bool bForce)
        {
            _ForceLog = bForce;
        }

        public bool IsConsolePrintingEnabled()
        {
            return LogCenter.Get().IsConsolePrintingEnabled(this.m_name);
        }

        public bool IsFilePrintingEnabled()
        {
            return LogCenter.Get().IsFilePrintingEnabled(this.m_name);
        }

        public bool IsPrintingEnabled()
        {
            return LogCenter.Get().IsPrintingEnabled(this.m_name);
        }

        public bool IsScreenPrintingEnabled()
        {
            return LogCenter.Get().IsScreenPrintingEnabled(this.m_name);
        }

        public void LogLevelPrint(string message, LogLevel logLevel)
        {
            LogCenter.Get().Print(this.m_name, message, logLevel);
        }

        bool NeedPrint()
        {
            return _ForceLog || LogCenter.Get().IsNeedPrint(this.m_name);
        }

        public void Print(LogLevel logLevel, string format, params object[] args)
        {
            if (!NeedPrint())
                return;
            if (!LogCenter.Get().CanLog(this.m_name, logLevel))
                return;
            string message = safeFormat(format, args);
            LogCenter.Get().Print(this.m_name, message, logLevel);
        }

        private string safeFormat(string format, params object[] args)
        {
            try
            {
                if(args.Length == 0)
                {
                    return format;
                }
                return string.Format(format, args);
            }
            catch(Exception e)
            {
                return "bad format:" + format + " exception:" + e;
            }            
        }

        public void ResetFile()
        {
            LogCenter.Get().ResetFile(this.m_name);
        }

        public void Debug(string format, params object[] args)
        {
            this.Print(LogLevel.Debug, format, args);
        }

        public void Info(string format, params object[] args)
        {
            this.Print(LogLevel.Info, format, args);
        }

        public void Warning(string format, params object[] args)
        {
            this.Print(LogLevel.Warning, format, args);
        }

        public void Error(string format, params object[] args)
        {
            this.Print(LogLevel.Error, format, args);
        }        

        public void ScreenPrint(string format, params object[] args)
        {
            string message = string.Format(format, args);
            LogCenter.Get().ScreenPrint(this.m_name, message, LogLevel.Debug);
        }
    }
}

