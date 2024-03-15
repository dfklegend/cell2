using System;

namespace Phoenix.Core
{
    [System.Reflection.ObfuscationAttribute(Exclude = true)]
    public class Singleton<T> where T : class
    {
        private static readonly T _instance;

        static Singleton()
        {            
            _instance = (T)Activator.CreateInstance(typeof(T), true);
        }

        public static T It
        {
            get
            {
                return Singleton<T>._instance;
            }
        }
    }
}

