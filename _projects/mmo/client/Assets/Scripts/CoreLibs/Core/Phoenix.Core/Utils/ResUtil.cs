using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;

namespace Phoenix.Core
{	
    public static class ResUtil
    {
        static string rawLoadTextFile(string strPath)
        {
            StreamReader sr = new StreamReader(strPath, System.Text.Encoding.UTF8);
            string content = sr.ReadToEnd();
            sr.Close();
            return content;
        }

        public static string LoadTextFile(string strPath)
        {
            if (strPath.Length == 0)
                return string.Empty;
            try
            {
                return rawLoadTextFile(strPath);
            }
            catch (System.Exception e)
            {
                PConsole.Log(e);
                return string.Empty;
            }
        }

        static void rawWriteTextFile(string path, string content)
        {
            StreamWriter sw = new StreamWriter(path, false, System.Text.Encoding.UTF8);
            sw.Write(content);
            sw.Close();
        }

        public static void WriteTextFile(string path, string content)
        {
            if (string.IsNullOrEmpty(path))
                return;
            try
            {
                rawWriteTextFile(path, content);
            }
            catch (System.Exception e)
            {
                PConsole.Log(e);
            }
        }
    }
} // namespace Phoenix
