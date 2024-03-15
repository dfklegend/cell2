namespace Phoenix.Core
{
    /**
     * @class HEvent
     * 
     * @brief The Event class is used as the base class for the creation of Event objects,
     * which are passed as parameters to event listeners when an event occurs.
     * 使用HEvent扩展定义能更明确的约定事件参数
     */
    public class HEvent<T>
    {
        private T _eventType;

        /// Constructor. Initialize <eventName_> to @a <eventName>.
        public HEvent(T eventType)
        {
            this._eventType = eventType;
        }

        /// Get name of the event.
        public T EventType
        {
            get { return this._eventType; }
        }
    }

    /**
     * @class HEventWithParams
     * 
     * @brief 
     * 可变参数
     */
    public class HEventWithParams<T> : HEvent<T>
    {      
        private object[] _params;
		public object[] paras
		{
			get
			{
				return _params;
			}
		}
        /// Constructor. Initialize <eventName_> to @a <eventName>.
        public HEventWithParams(T eventType, params object[] args)
            :base(eventType)
        {
            _params = args;
        }        
    }


    public static class HEventUtil
    {
        public static void Dispatch<T>(EventCenter<T> center, HEvent<T> hEvent)
        {
            center.Dispatch(hEvent.EventType, hEvent);
        }
    }
}

