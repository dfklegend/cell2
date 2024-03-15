using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using Phoenix.Core;
using UnityEngine;

namespace Phoenix.Client
{
    public class UnityConsole : IConsole
    {
        public void Log(object message)
        {
            Debug.Log(message);
        }

        public void Warning(object message)
        {
            Debug.LogWarning(message);
        }

        public void Error(object message)
        {
            Debug.LogError(message);
        }
    }
}
