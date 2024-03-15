using System;

namespace Phoenix.Utils
{

    public static class RandUtil
    {
        private static Random _random = new Random();

        public static double Range(double min, double max)
        {
            var v = _random.NextDouble();
            return min + v * (max - min);
        }

        public static int Range(int min, int max)
        {
            return _random.Next(min, max);
        }
    }
} // Phoenix.Utils