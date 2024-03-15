using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator;

namespace Phoenix.Game
{
    public class UIUpdater : AppTask.BaseCommonTask
    {
        public override void Update()
        {
            UIMgr.It.Update();
        }

        public override void LateUpdate()
        {
            UIMgr.It.LateUpdate();
        }
    }
} // namespace Phoenix
