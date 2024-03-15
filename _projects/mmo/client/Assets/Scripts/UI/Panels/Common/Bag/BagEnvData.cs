using System.Collections.Generic;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator.BagSystem;
using Phoenix.Core;

namespace Phoenix.Game
{
    public class BagEnvData
    {
        // 选择类型
        public int selectMode = 0;
        // 窗口位置
        // 构建的数据
        public BagViewData viewData;

        // 当前item
        public IShowItem curItem;

        // 处理器
        public IShowItemHandler handler;
    }     
} // namespace Phoenix
