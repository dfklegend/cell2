using System.Collections.Generic;
using System;

namespace Phoenix.Core
{ 
    public delegate void EventDelegate(params object[] args);
    
    interface IEventListenerList
    {

        void Add(EventDelegate handler);
        void Remove(EventDelegate handler);
        bool IsHandlersNull();
        void Handle(params object[] args);
    }
    class EventListenerList : IEventListenerList     
    {
        private EventDelegate _handlers;
        public event EventDelegate handlers {
            add
            {
                if (_handlers == null)
                {
                    _handlers = value;
                }
                else
                {
                    // 避免重复添加
                    var list = _handlers.GetInvocationList();
                    for (var i = 0; i < list.Length; i ++)
                    {
                        if (list[i].Equals(value))
                            return;
                    }
                    _handlers += value;                    
                }
            }
            remove
            {
                _handlers -= value;
            }
        }
        public bool IsHandlersNull()
        {
            return _handlers == null;
        }
        public void Handle(params object[] args)
        {
            try
            {
                // 这里是否会有 处理过程中自删的问题?
                if (_handlers != null)
                {
                    //Debug.Log("handlers!");
                    _handlers(args);
                    //Debug.Log("handlers end!");
                }
            }
            catch(System.Exception ex )
            {
                //Debug.LogException( ex );
                PConsole.Error(ex);
            }
        }

        public void Add(EventDelegate handler)
        {
            // 避免重复添加
            if (_handlers == null)
            {
                _handlers = handler;
                return;
            }
                
            var list = _handlers.GetInvocationList();
            for (var i = 0; i < list.Length; i++)
            {
                if (list[i].Equals(handler))
                    return;
            }
            _handlers += handler;
        }

        public void Remove(EventDelegate handler)
        {
            _handlers -= handler;
        }

        public void LogHandlers()
        {
            
        }
    }

    class OneHandler
    {
        public EventDelegate handler;
        public bool removed = false;

        public OneHandler(EventDelegate h)
        {
            handler = h;
        }
    }

    // 上面直接用event,事件处理过程中，移除并不会马上移除
    // 保证触发过程中 移除，注册不会丢掉触发
    class EventListenerListEx : IEventListenerList
    {
        private List<OneHandler> _handlers = new List<OneHandler>();
        private bool _handling = false;
        private bool _hasHandleRemoved = false;
        
        public bool IsHandlersNull()
        {
            return _handlers.Count == 0;
        }
        public void Handle(params object[] args)
        {
            if (_handlers.Count == 0)
                return;
            _handling = true;
            for(var i = 0; i < _handlers.Count; i ++)
            {
                handleOne(_handlers[i], args);
            }
            _handling = false;

            applyRemoved();
        }               

        private void handleOne(OneHandler one, params object[] args)
        {
            try 
            {
                if(!one.removed)
                    one.handler(args);
            }
            catch(Exception e)
            {
                PConsole.Error(e);
            }
        }

        private void applyRemoved()
        {
            if (!_hasHandleRemoved)
                return;
            _hasHandleRemoved = false;

            // remove removed
            for (var i = 0; i < _handlers.Count;)
            {
                var one = _handlers[i];
                if (one.removed)
                    _handlers.RemoveAt(i);
                else
                    i++;
            }
        }

        public void LogHandlers()
        {
        }
        
        int FindIndex(EventDelegate handler)
        {
            for(var i = 0; i < _handlers.Count; i ++)
            {
                if (_handlers[i].handler.Equals(handler))
                    return i;
            }
            return -1;
        }
        
        public void Add(EventDelegate handler)
        {
            var index = FindIndex(handler);
            if (index != -1)
            {
                var one = _handlers[index];
                if (one.removed)
                    one.removed = false;
                return;
            }
                
            _handlers.Add(new OneHandler(handler));
        }

        public void Remove(EventDelegate handler)
        {
            var index = FindIndex(handler);
            if ( index == -1)
                return;
            if(_handling)
            {
                _handlers[index].removed = true;
                _hasHandleRemoved = true;
                return;
            }

            _handlers.RemoveAt(index);
        }
    }

    public class EventCenter<TKey>
    {
        Dictionary<TKey, IEventListenerList> _dictMsgMap = new Dictionary<TKey, IEventListenerList>();

        IEventListenerList GetHandlerList(TKey type)
        {
            try
            {
                if (_dictMsgMap.ContainsKey(type))
                    return _dictMsgMap[type];
                else
                {
                    _dictMsgMap[type] = createList();
                    return _dictMsgMap[type];
                }
            }
            catch (System.Exception /*ex*/)
            {
                return null;
            }
        }

        private IEventListenerList createList()
        {
            return new EventListenerListEx();
        }

        public void Bind(TKey type, EventDelegate handler, bool bAddHandler = true)
        {
            //Debug.Log("bind:" + type + " hanler:" + handler );
            IEventListenerList l = GetHandlerList(type);
            if (bAddHandler)
            {
                l.Add(handler);
            }                
            else
            {
                //l.handlers -= handler;
                l.Remove(handler);
            }                
        }

        public bool Dispatch(TKey type, params object[] args)
        {
            if (!_dictMsgMap.ContainsKey(type))
            {
                //Debug.Log("can not find key:" + type);
                return false;
            }
            IEventListenerList l = GetHandlerList(type);
            if (l.IsHandlersNull())
                return false;
            l.Handle(args);
            return true;
        }        
    }
}