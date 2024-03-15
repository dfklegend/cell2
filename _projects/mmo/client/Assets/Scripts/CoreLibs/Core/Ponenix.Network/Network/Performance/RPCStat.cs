using System;
using System.Text;
using TimeUtil = Phoenix.Utils.TimeUtil;
using Phoenix.Core;

namespace Phoenix.Network
{
    // 功能说明: 
    // 客户端
    // 统计一下RPC调用
    // 平均每秒多少RPC调用
    // 平均响应时间
    public class RPCStat
    {
        public class Stat
        {
            public float totalCost;
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

        public void AddOneCost(float cost)
        {
            _total.callTimes++;
            _total.totalCost += cost;

            _phase.callTimes++;
            _phase.totalCost += cost;
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
            var now = TimeUtil.Now();
            if (now < _nextDumpTime)
                return;
            _nextDumpTime = now + DUMP_TIME;
            PConsole.Log(makeStat(_lastPhase, PHASE_TIME));
        }

        private string makeStat(Stat stat, float totalTime)
        {
            if (stat.callTimes == 0)
                return "null state, no rpc called";
            StringBuilder sb = new StringBuilder();
            sb.Append($"  called: {stat.callTimes}\r\n");
            sb.Append($"  avg response time: {stat.totalCost/stat.callTimes}\r\n");
            sb.Append($"  speed: {stat.callTimes/ totalTime}(calls/s)\r\n");
            return sb.ToString();
        }

        public void DumpTotal()
        {
            PConsole.Log("---- total stat ----");
            var now = TimeUtil.GetSystemSecond();
            PConsole.Log(makeStat(_total, now - _startTime));
        }

        public int GetTotalTimes()
        {
            return _total.callTimes;
        }
    }
}

