namespace Phoenix.Core
{
    // 提供获取当前RunEnv支持
    public interface IRunEnvGetter
    {
        RunEnv Get();
    }
    
    // 单Env
    public class OneEnvAppGetter : IRunEnvGetter
    {
        private RunEnv _env;
        public OneEnvAppGetter(RunEnv env)
        {
            _env = env;
        }

        public RunEnv Get()
        {
            return _env;
        }
    }

    // 基于线程对应的Env
    // 假设未来服务器需要 多逻辑线程逻辑
    public class ThreadEnvGetter : IRunEnvGetter
    {   
        public ThreadEnvGetter()
        {
            
        }

        public RunEnv Get()
        {
            // . 获取当前线程id
            // . 获取对应的Env
            return null;
        }
    }
}
