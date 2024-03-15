using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using System.Net.Sockets;
using Network;


namespace Phoenix.Game.Card
{
    public class HEventCharInfo: HEvent<int>
    {
        public Cproto.CharInfo charInfo;
        public HEventCharInfo(Cproto.CharInfo info)
            : base(EventDefine.CharInfo)
        {
            charInfo = info;
        }
    }

    public class HEventBattleLog : HEvent<int>
    {
        public string log;
        public HEventBattleLog(string info)
            : base(EventDefine.BattleLog)
        {
            log = info;
        }
    }

    public class HEventRefreshSceneInfo : HEvent<int>
    {        
        public HEventRefreshSceneInfo()
            : base(EventDefine.RefreshSceneInfo)
        {            
        }
    }

}// namespace Phoenix
