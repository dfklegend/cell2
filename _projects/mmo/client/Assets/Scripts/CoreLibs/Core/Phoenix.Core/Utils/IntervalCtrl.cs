namespace Phoenix.Utils
{

    // 控制间隔
    // 考虑整数范围的问题
    // 配合Environment.TickCount使用
    public class IntervalCtrl
    {
        private int _lastDo = 0;

        public bool CanDo(int now, int interval)
        {            
            // 考虑tick 回转
            if (now < _lastDo + interval && now > _lastDo)
                return false;
            _lastDo = now;
            return true;
        }
    }

    public class LongIntervalCtrl
    {
        private long _lastDo = 0;

        public bool CanDo(long now, long interval)
        {
            // 考虑tick 回转
            if (now < _lastDo + interval && now > _lastDo)
                return false;
            _lastDo = now;
            return true;
        }
    }
} // Phoenix.Utils