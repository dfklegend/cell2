namespace Phoenix.Core
{
    public static class AppEnv
    {
        private static string _rootDir = "./";
        public static string rootDir { get { return _rootDir; } }

        private static IRunEnvGetter _runEnvGetter;
        private static IUniqueIDAllocer _uniqueIdAllocer =
            new SimpleUniqueIDAllocer();
        
        
        public static void SetRootDir(string dir)
        {
            dir = dir.TrimEnd('/', '\\');
            _rootDir = dir;
        }

        public static string GetDataPath()
        {
            return _rootDir;
        }

        public static string GetFullNameInDataPath(string path)
        {
            return $"{_rootDir}/{path}";
        }

        public static void InitRunEnvGetter(IRunEnvGetter getter)
        {
            _runEnvGetter = getter;
        }

        public static RunEnv GetRunEnv()
        {
            return _runEnvGetter.Get();
        }

        public static void InitUniqueIDAllocer(IUniqueIDAllocer allocer)
        {
            _uniqueIdAllocer = allocer;
        }

        public static IUniqueIDAllocer GetAllocer()
        {
            return _uniqueIdAllocer;
        }

        public static ulong AllocUniqueId()
        {
            return _uniqueIdAllocer.AllocId();
        }
    }
}
