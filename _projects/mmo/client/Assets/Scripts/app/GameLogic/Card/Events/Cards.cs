using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using System.Net.Sockets;
using Network;


namespace Phoenix.Game.Card
{       
    public class HEventServerRefreshCards : HEventServerMsg<Cproto.RefreshCards>
    {
        public HEventServerRefreshCards(Cproto.RefreshCards msg)
            : base(msg, EventDefine.RefreshServerCards)
        {
        }
    }

    public class HEventRefreshCards : HEvent<int>
    {
        public HEventRefreshCards()
            : base(EventDefine.RefreshCards)
        {
        }
    }

    public class HEventServerAttrsChanged : HEventServerMsg<Cproto.UnitAttrsChanged>
    {
        public HEventServerAttrsChanged(Cproto.UnitAttrsChanged msg)
            : base(msg, EventDefine.UnitAttrsChanged)
        {
        }
    }
}// namespace Phoenix
