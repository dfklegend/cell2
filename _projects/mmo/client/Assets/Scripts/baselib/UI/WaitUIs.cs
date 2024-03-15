using System.Collections.Generic;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Core;

namespace Phoenix.Game
{
    // 等待一系列的界面载入完毕
    public class WaitUIs
    {
        private List<string> _waits = new List<string>();
        private TimerID _timer;

        public void Add(string name)
        {
            _waits.Add(name);
        }

        public void StartWait(System.Action callback)
        {
            _timer = AppEnv.GetRunEnv().timer.AddTimer((args) =>
            {
                if (checkReady())
                {
                    AppEnv.GetRunEnv().timer.Cancel(_timer);
                    _timer = null;
                    callback?.Invoke();
                }
                    
            }, 0.001f, 0.001f);
        }

        private bool checkReady()
        {
            foreach(var one in _waits)
            {
                var panel = UIMgr.It.GetPanel(one);
                if (panel != null && !panel.IsReady())
                    return false;
            }
            return true;
        }
    }
} // namespace Phoenix
