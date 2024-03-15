using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game.FText
{
    public class ShowerFactory : Singleton<ShowerFactory>
    {
        private IntToClassFactory<FloatTextShower> _factory =
            new IntToClassFactory<FloatTextShower>();
        public FloatTextShower Create(int t, params object[] args)
        {
            FloatTextShower shower = _factory.Create(t);
            if (shower == null)
                return null;
            shower.Create(args);
            return shower;
        }
    }
} // namespace Phoenix
