namespace Phoenix.Core
{
    
    // 非线程安全
    public class GlobalEvents : Singleton<GlobalEvents>
    {
        private EventCenter<int> _events = new EventCenter<int>();
        public EventCenter<int> events { get { return _events; } }

        private EventCenter<string> _strEvents = new EventCenter<string>();
        public EventCenter<string> strEvents { get { return _strEvents; } }
    }
}

