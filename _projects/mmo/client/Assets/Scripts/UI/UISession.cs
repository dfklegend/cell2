using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator;

namespace Phoenix.Game
{
    public class UISession : Singleton<UISession>
    {
        public eEquipSlot curSlot = eEquipSlot.MainHand;
    }
} // namespace Phoenix
