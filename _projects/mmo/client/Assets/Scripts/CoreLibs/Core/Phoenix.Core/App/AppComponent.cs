using System;
using System.Collections.Generic;
using Phoenix.Utils;

namespace Phoenix.Core
{
    public class AppComponentPriority
    {
        public const int BelowNormal = 50;
        public const int Normal = 100;
        public const int AboveNormal = 150;
        public const int Highest = 9999;
    }

    // 游戏启动时，注册并按按顺序初始化
    // 组合AppComponent来组合成不同App
    // 按优先级依次start
    // 前一个完成，才开始后面的
    // Stop反向关闭顺序

    public interface IAppComponent
    {
        int GetPriority();

        void Start();
        bool IsReady();

        // 准备期的update
        void PrepareUpdate();
        void Update();
        void StopUpdate();
        void Stop();
        // 是否停止完毕
        bool IsStopped();
    }
    
    public static class AppComponentMgr
    {
        private static List<IAppComponent> _prepares = new List<IAppComponent>();
        private static IAppComponent _waitReady;

        private static bool _modulesDirt = false;
        private static List<IAppComponent> _modules = new List<IAppComponent>();

        private static bool _stop = false;
        private static IAppComponent _waitStop;

        public static void Clear()            
        {
            _waitReady = null;
            _waitStop = null;
            _prepares.Clear();
            _modules.Clear();
        }

        public static void Reset()
        {           
            _modules.Clear();
            _modulesDirt = false;
            _stop = false;
        }


        public static void Register<T>()
            where T: class, IAppComponent
        {
            var obj = Activator.CreateInstance(typeof(T)) as IAppComponent;
            _prepares.Add(obj);
        }

        private static void sort(List<IAppComponent> modules)
        {
            modules.Sort((a, b) => b.GetPriority() - a.GetPriority());
        }

        public static void StartAll()
        {          
            sort(_prepares);
            TryStart();
        }

        private static bool hasModuleWaiting()
        {
            return _waitReady != null;
        }

        private static bool isStarting()
        {
            return _prepares.Count > 0 || hasModuleWaiting();
        }

        private static void addToReadys(IAppComponent module)
        {
            _modules.Add(module);
            _modulesDirt = true;
        }

        private static void TryStart()
        {
            while (startOne());            
        }

        // return
        //     true: instant task
        private static bool startOne()
        {
            if (_prepares.Count == 0)
                return false;
            if (hasModuleWaiting())
                return false;
            var head = _prepares[0];
            _prepares.RemoveAt(0);

            head.Start();
            if (head.IsReady())
            {
                addToReadys(head);
                return true;
            }

            // wait
            _waitReady = head;
            return false;
        }

        private static void checkWaiting()
        {
            if (_waitReady == null)
                return;
            if(_waitReady.IsReady())
            {   
                addToReadys(_waitReady);

                _waitReady = null;
            }
            else 
            {
                _waitReady.PrepareUpdate();
            }
        }

        public static bool IsAllReady()
        {
            if (_prepares.Count > 0)
                return false;
            if (_waitReady != null)
                return false;
            return true;
        }        

        public static void Update()
        {
            if(_stop)
            {
                updateStop();
                return;
            }
            checkWaiting();
            TryStart();            
            updateReadys();
        }

        private static void updateReadys()
        {
            if (_modulesDirt)
            {
                sort(_modules);
                _modulesDirt = false;
            }
                
            for (var i = 0; i < _modules.Count; i++)
            {
                _modules[i].Update();
            }
        }

        public static void StopAll()
        {            
            if(isStarting())
            {
                PConsole.Error("Error, Stop while starting!");
                return;
            }
            _stop = true;
            // 依次stop，从后到前
            stopLast();
        }

        private static void updateStop()
        {
            if(_waitStop != null)
            {
                _waitStop.StopUpdate();
                if(_waitStop.IsStopped())
                {
                    _waitStop = null;
                    stopLast();
                }
            }
        }

        private static void stopLast()
        {
            if (_modules.Count == 0)
                return;            
            _waitStop = _modules[_modules.Count - 1];
            _modules.RemoveAt(_modules.Count - 1);

            _waitStop.Stop();
        }

        public static bool IsAllStopped()
        {
            return _waitStop == null && _modules.Count == 0;
        }
    }

} // namespace Phoenix.Module