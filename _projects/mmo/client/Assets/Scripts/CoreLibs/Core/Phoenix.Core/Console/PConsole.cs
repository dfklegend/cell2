using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{
    // 基本输出信息
    public static class PConsole
    {
        private static IConsole _console = new DummyConsole();
        public static void SetConsole(IConsole console) { _console = console; }
        public static void Log(object message)
        {
            _console.Log(message);
        }

        public static void Warning(object message)
        {
            _console.Warning(message);
        }

        public static void Error(object message)
        {
            _console.Error(message);
        }
    }
}
