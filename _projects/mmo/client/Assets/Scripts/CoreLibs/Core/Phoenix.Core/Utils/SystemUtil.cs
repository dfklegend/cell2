using System;
using System.Collections;
using System.Collections.Generic;
using System.IO;
using System.Reflection;
using System.Text;

namespace Phoenix.Utils
{
    public static class SystemUtil
    {
        // 获取所有某个类的子类
        public static List<Type> GetAllClass<T>() where T : class
        {
            List<Type> ret = new List<Type>();
            foreach (Assembly assembly in AppDomain.CurrentDomain.GetAssemblies())
            {
                try
                {
                    if (assembly.IsDynamic)
                    {
                        //Debug.Log("dynamic assembly:" + assembly.FullName);
                        continue;
                    }
                        
                    foreach (System.Type type in assembly.GetExportedTypes())
                    {
                        if (((typeof(T).IsAssignableFrom(type) && type.IsClass) && !type.IsAbstract))
                        {
                            ret.Add(type);
                        }
                    }
                }
                catch(Exception e) 
                {
                    Core.PConsole.Error(e);
                }
                           
            }
            return ret;
        }

        public static void LogHandledException(Exception e)
        {            
            Log.LogCenter.Exception.Error("HandledException:");
            Log.LogCenter.Exception.Error(e.ToString());
        }

        public static void AddUnhandledExceptionLog()
        {
            AppDomain.CurrentDomain.UnhandledException += CurrentDomain_UnhandledException;
        }

        private static void CurrentDomain_UnhandledException(object sender, UnhandledExceptionEventArgs e)
        {
            Log.LogCenter.Exception.Error("---- UnhandledException: ----");
            Log.LogCenter.Exception.Error(e.ExceptionObject.ToString());
        }

        public static void ChangeLogFile(string logType, string name)
        {
            Log.LogCenter.GetLogInfo(logType, true).SetFileName(name);
        }
    }
}
