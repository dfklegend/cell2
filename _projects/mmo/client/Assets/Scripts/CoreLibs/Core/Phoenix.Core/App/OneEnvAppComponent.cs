using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{
    // 初始化成逻辑单线程
    public class OneEnvAppComponent : IAppComponent
    {
        RunEnv _env;
        public RunEnv env { get { return _env; } }        

        public int GetPriority()
        {
            return AppComponentPriority.Highest;
        }

        public void Start()
        {
            _env = new RunEnv();
            _env.CreateEnv();

            AppEnv.InitRunEnvGetter(new OneEnvAppGetter(_env));
        }
        public bool IsReady()
        {
            return true;
        }

        // 准备期的update
        public void PrepareUpdate() { }
        public void Update() 
        {
            _env.Update();
        }

        public void StopUpdate()
        {

        }

        public void Stop() { }
        // 是否停止完毕
        public bool IsStopped()
        {
            return true;
        }
    }
}
