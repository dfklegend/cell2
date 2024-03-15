using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{
    // 基本输出信息
    public interface IConsole
    {
        void Log(object message);
        void Warning(object message);
        void Error(object message);
    }

    public class DummyConsole : IConsole
    {
        public void Log(object message) 
        {
            Console.WriteLine(message);
        }

        public void Warning(object message) 
        {
            Console.WriteLine(message);
        }

        public void Error(object message) 
        {
            Console.WriteLine(message);
        }
    }
}
