using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{
    // 定义一些需要初始化的接口
    // 汇总在一起
    public static class AppInit
    {
        public static void SetConsole(IConsole console)
        {
            PConsole.SetConsole(console);
        }

        public static void SetAppDir(string dir)
        {
            AppEnv.SetRootDir(dir);
        }
    }
}
