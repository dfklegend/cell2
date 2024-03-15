using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Utils;


namespace Phoenix.Game.Card.FightStat
{
    public class Consts
    {
        // 1秒抽样一次
        public const int SampleInterval = 1000;
    }

    public class UnitStat
    {
        public string name = "";
        public int dmg = 0;
    }

    public class SampleData
    {
        public int dmg = 0;
    }

    public class UnitSample
    {
        public List<SampleData> data = new List<SampleData>();
    }


    // 伤害统计    
    public class FightStat
    {
        int _startTime;
        int _nextSampleTime;
        int _sampleIndex = 0;
        bool _stopped;

        // serverId:
        private Dictionary<int, UnitStat> _units = new Dictionary<int, UnitStat>();
        private Dictionary<int, UnitSample> _samples = new Dictionary<int, UnitSample>();

        public Dictionary<int, UnitStat> GetUnits()
        {
            return _units;
        }

        public int GetUnitNum()
        {
            return _units.Count;
        }

        public Dictionary<int, UnitSample> GetSamples()
        {
            return _samples;
        }

        public int GetSampleNum()
        {
            return _sampleIndex;
        }


        public void Start()
        {
            _startTime = TimeUtil.NowTick();
            _nextSampleTime = 0;
            _sampleIndex = 0;
            _units.Clear();
            _samples.Clear();
            _stopped = false;
        }

        private int now()
        {
            return TimeUtil.NowTick() - _startTime;
        }            

        public void Stop()
        {
            _stopped = true;
        }

        public void AddUnitInfo(int id, string name)
        {
            UnitStat unit;
            if (_units.TryGetValue(id, out unit))
            {
                unit.name = name;
            }
            else
            {
                unit = new UnitStat();
                unit.name = name;
                _units[id] = unit;
                onNewUnit(id);
            }
        }

        public void AddDmg(int id, int dmg)
        {
            UnitStat unit;
            if(_units.TryGetValue(id, out unit))
            {
                unit.dmg += dmg;
            }
            else
            {
                unit = new UnitStat();
                unit.dmg = dmg;
                _units[id] = unit;
                onNewUnit(id);
            }
        }

        public void Update()
        {
            if (_stopped)
                return;
            tryMakeSample();
        }

        public void tryMakeSample()
        {
            var now = this.now();
            if(now < _nextSampleTime)
            {
                return;
            }

            _nextSampleTime += Consts.SampleInterval;
            makeSample();
            _sampleIndex++;
        }

        // 尝试抽样一次
        public void makeSample()
        {
            foreach(KeyValuePair<int, UnitStat> kv in _units)
            {
                var sample = _samples[kv.Key];

                sample.data.Add(new SampleData()
                {
                    dmg = kv.Value.dmg
                });
            }
        }

        private void onNewUnit(int id)
        {            
            var sample = new UnitSample();
            _samples[id] = sample;

            // 补充一下samples
            if (_sampleIndex == 0)
                return;
            for (var i = 0; i < _sampleIndex; i++) 
            {
                sample.data.Add(new SampleData());
            }
        }
    }
        
} // namespace Phoenix
