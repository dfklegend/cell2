using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using Phoenix.Utils;

namespace Phoenix.Log
{
    public static class LogEnv
    {
        private static string _logPrefix = "";
        private static string _rootPath = "";

        public static void SetLogPrefix(string prefix)
        {
            _logPrefix = prefix;
        }

        public static string MakeLogFileName(string name)
        {
            if (string.IsNullOrEmpty(_logPrefix))
                return name;
            return $"{_logPrefix}.{name}";
        }

        public static void SetRootPath(string path)
        {
            _rootPath = path;
        }

        public static string GetDataPath()
        {
            return _rootPath;
        }

        public static string GetFullNameInDataPath(string path)
        {
            return $"{_rootPath}/{path}";
        }
    }
}