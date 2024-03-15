using System;
using System.Text;
using TimeUtil = Phoenix.Utils.TimeUtil;

namespace Phoenix.Network
{   
    // 服务器
    // 统计一下处理了多少RPC
    // 平均每秒多少RPC调用    
    public class ServerStat
    {
        public class Stat
        {
            public int callTimes;
        }

        const float PHASE_TIME = 5.0f;
        const float DUMP_TIME = 30f;

        private float _startTime;
        private Stat _total;
        private float _phaseStartTime;
        private Stat _phase;
        private Stat _lastPhase;
        private float _nextDumpTime = 0f;

        public void Start()
        {
            _startTime = TimeUtil.GetSystemSecond();
            _phaseStartTime = _startTime;
            _nextDumpTime = _startTime + DUMP_TIME;
            _total = new Stat();
            _phase = new Stat();
            _lastPhase = new Stat();
        }

        public void AddOne()
        {
            _total.callTimes++;           
            _phase.callTimes++;            
            tryStartNewPhase();
        }

        private void tryStartNewPhase()
        {
            var now = TimeUtil.GetSystemSecond();
            if (now < _phaseStartTime + PHASE_TIME)
                return;

            _phaseStartTime = now;
            _lastPhase = _phase;
            _phase = new Stat();
        }

        public void TryDump()
        {
            var now = TimeUtil.GetSystemSecond();
            if (now < _nextDumpTime)
                return;
            _nextDumpTime = now + DUMP_TIME;

            var result = makeStat(_lastPhase, PHASE_TIME);
            if (!string.IsNullOrEmpty(result))
            {
                Log.LogCenter.Default.Info(result);
            }
        }

        private string makeStat(Stat stat, float totalTime)
        {
            if (stat.callTimes == 0)
                return "";
            StringBuilder sb = new StringBuilder();
            sb.Append($"  called: {stat.callTimes}\r\n");            
            sb.Append($"  speed: {stat.callTimes / totalTime}(calls/s)\r\n");
            return sb.ToString();
        }

        public void DumpTotal()
        {
            Console.WriteLine("---- total stat ----");
            var now = TimeUtil.GetSystemSecond();
            Console.WriteLine(makeStat(_total, now - _startTime));
        }
    }
}

