using System;
using System.IO;
namespace Phoenix.Log
{
    public class LogInfo
    {
        private bool m_consolePrinting;
        private bool m_fileCreateFailed;
        private bool m_filePrinting;
        private LogLevel m_logLevel;
        private string m_name;
        private string m_fileName;
        private bool m_screenPrinting;
        private StreamWriter m_writer;

        public LogInfo(string name)
        {
            this.m_logLevel = LogLevel.Debug;
            this.m_name = name;
            initFileName(name);
        }

        public LogInfo(string name, bool consolePrinting, bool screenPrinting, bool filePrinting, LogLevel logLevel)
        {
            this.m_logLevel = LogLevel.Debug;
            this.m_name = name;
            initFileName(name);
            this.m_consolePrinting = consolePrinting;
            this.m_screenPrinting = screenPrinting;
            this.m_filePrinting = filePrinting;
            this.m_logLevel = logLevel;
        }

        public bool DidFileCreateFail()
        {
            return this.m_fileCreateFailed;
        }

        public StreamWriter GetFileWriter()
        {
            return this.m_writer;
        }

        public LogLevel GetLogLevel()
        {
            return this.m_logLevel;
        }

        public string GetName()
        {
            return this.m_name;
        }

        private void initFileName(string name)
        {
            m_fileName = LogEnv.MakeLogFileName(name);
        }

        public void SetFileName(string n)
        {
            m_fileName = n;
            TryCloseFileWriter();
        }

        public void UpdateFileName()
        {
            initFileName(m_name);
            TryCloseFileWriter();
        }

        public string GetFileName()
        {
            return m_fileName;
        }

        public void TryCloseFileWriter()
        {
            if (this.m_writer != null)
            {
                m_writer.Close();
                m_writer = null;
            }
        }

        public bool IsConsolePrintingEnabled()
        {
            return this.m_consolePrinting;
        }

        public bool IsFilePrintingEnabled()
        {
            //return (ApplicationMgr.IsInternal() && this.m_filePrinting);
            return (this.m_filePrinting);
        }

        public bool IsPrintingEnabled()
        {
            return (this.IsConsolePrintingEnabled() || this.IsScreenPrintingEnabled());
        }

        public bool IsScreenPrintingEnabled()
        {
            return this.m_screenPrinting;
        }

        public LogInfo SetConsolePrintingEnabled(bool enable)
        {
            this.m_consolePrinting = enable;
            return this;
        }

        public void SetFileCreateFailed()
        {
            this.m_fileCreateFailed = true;
        }

        public LogInfo SetFilePrintingEnabled(bool enable)
        {
            this.m_filePrinting = enable;
            return this;
        }

        public void SetFileWriter(StreamWriter writer)
        {
            this.m_writer = writer;
        }

        public LogInfo SetLogLevel(LogLevel logLevel)
        {
            this.m_logLevel = logLevel;
            return this;
        }

        public LogInfo SetScreenPrintingEnabled(bool enable)
        {
            this.m_screenPrinting = enable;
            return this;
        }
    }
}

