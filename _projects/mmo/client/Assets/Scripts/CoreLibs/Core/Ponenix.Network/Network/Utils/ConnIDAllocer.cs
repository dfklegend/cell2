using System.Threading;

namespace Phoenix.Network
{
    public static class ConnIDAllocer
    {
        static long _nextId = 1;

        public static long Alloc()
        {
            return Interlocked.Increment(ref _nextId);            
        }
    }
}
