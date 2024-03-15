using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator.BagSystem;
using System;

namespace Phoenix.Game
{   
    // 某些情况下，比如不显示名字等
    // 可以通过此配置调整
    public interface IIconStyleOptions
    {
        bool NeedHideName();
    }

    public class DefaultOptions : IIconStyleOptions
    {  
        public bool NeedHideName()
        {
            return false;
        }
    }

    public class TipOptions: IIconStyleOptions
    {
        public bool NeedHideName()
        {
            return true;
        }
    }

    public static class IconStyleOptions
    {
        public static DefaultOptions Default = new DefaultOptions();
        public static TipOptions Tip = new TipOptions();
    }
} // namespace Phoenix
