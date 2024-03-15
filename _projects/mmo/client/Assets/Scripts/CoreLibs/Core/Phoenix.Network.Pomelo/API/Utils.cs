using System;
using System.Collections.Generic;
using System.Reflection;

namespace Phoenix.API
{
    public static class APIUtils
    {
        public static void LogException(Exception e)
        {
            //Console.WriteLine("Exception:");
            //Console.WriteLine(e);

            Log.LogCenter.Exception.Error(e.ToString());
        }

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
                catch (Exception e)
                {
                    LogException(e);
                }

            }
            return ret;
        }
    }
}

