using System;
using System.Diagnostics;

namespace Phoenix.Utils
{

    public static class TimeUtil
    {
        private static long _tickPerSecond = 1;
        private static long _tickPerMs = 1;
        static TimeUtil()
        {
            _tickPerSecond = Stopwatch.Frequency;
            _tickPerMs = _tickPerSecond / 1000;
        }

        // 获取系统second时间
        // 精度不高
        public static float GetSystemSecond()
        {
            return (float)Environment.TickCount / 1000f;
        }

        public static float Now()
        {
            return GetSystemSecond();
        }

        public static int NowTick()
        {
            return Environment.TickCount;
        }

        public static long HiNowMs()
        {
            return Stopwatch.GetTimestamp() / _tickPerMs;
        }

        public static float HiNow()
        {
            return (float)Stopwatch.GetTimestamp() / _tickPerSecond;
        }
    }
} // Phoenix.Utils