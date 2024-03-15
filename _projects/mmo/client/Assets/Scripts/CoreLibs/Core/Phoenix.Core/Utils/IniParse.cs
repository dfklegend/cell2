using System.Collections;
using System.IO;
using System.Text;
using Phoenix.Core;

namespace Phoenix.Utils
{
    // 提供ini文件的解析读取
    public class IniUtil
    {       
        static string rawLoadTextFile(string strPath)
        {
            StreamReader sr = new StreamReader(strPath);
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
                Core.PConsole.Error(e);
                return string.Empty;
            }
        }                

        public static bool ParseIniFile(string filePath, ConfigFileEntryParseCallback callback, object userData = null)
        {
            if (callback == null)
            {
                PConsole.Warning("FileUtils.ParseConfigFile() - no callback given");
                return false;
            }
            if (!System.IO.File.Exists(filePath))
            {
                PConsole.Warning(string.Format("FileUtils.ParseConfigFile() - file {0} does not exist", filePath));
                return false;
            }

            try
            {
                string content = LoadTextFile(filePath);
                bool bRet = ParseConfigText(content, callback, userData);                
                return bRet;
            }
            catch(System.Exception e)
            {
                PConsole.Error(e);
                return false;
            }            
        }

        public static bool ParseConfigText(string content, ConfigFileEntryParseCallback callback, object userData = null)
        {
            if (callback == null)
            {
                PConsole.Warning("FileUtils.ParseConfigFile() - no callback given");
                return false;
            }

            MemoryStream s = new MemoryStream(System.Text.Encoding.UTF8.GetBytes(content));
            bool bRet = ParseConfigStream(s, callback, userData);
            s.Close();
            return bRet;
        }

        public static bool ParseConfigStream(Stream s, ConfigFileEntryParseCallback callback, object userData = null)
        {
            if (callback == null)
            {
                PConsole.Warning("FileUtils.ParseConfigFile() - no callback given");
                return false;
            }           
            int num = 1;
            using (StreamReader reader =  new StreamReader(s))
            {
                string baseKey = string.Empty;
                while (reader.Peek() != -1)
                {
                    string str2 = reader.ReadLine().Trim();
                    if ((str2.Length >= 1) && (str2[0] != ';'))
                    {
                        if (str2[0] == '[')
                        {
                            if (str2[str2.Length - 1] != ']')
                            {
                                PConsole.Warning(string.Format("FileUtils.ParseConfigFile() - bad key name \"{0}\" on line {1}", str2, num));
                            }
                            else
                            {
                                baseKey = str2.Substring(1, str2.Length - 2);
                            }
                        }
                        else if (!str2.Contains("="))
                        {
                            // 非注释
                            if( str2.IndexOf( "//") == -1 )
                                PConsole.Warning(string.Format("FileUtils.ParseConfigFile() - bad value pair \"{0}\" on line {1}", str2, num));
                        }
                        else
                        {
                            char[] separator = new char[] { '=' };
                            string[] strArray = str2.Split(separator);
                            if (strArray.Length >= 2)
                            {
                                string subKey = strArray[0].Trim();
                                string val = strArray[1].Trim();
                                if ( val.Length > 0 && (val[0] == '"') && (val[val.Length - 1] == '"'))
                                {
                                    val = val.Substring(1, val.Length - 2);
                                }
                                callback(baseKey, subKey, val, userData);
                            }
                        }
                    }
                }
            }
            return true;
        }        

        public delegate void ConfigFileEntryParseCallback(string baseKey, string subKey, string val, object userData);
    }
}
