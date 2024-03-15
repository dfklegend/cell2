using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using Network;
using pbc = global::Google.Protobuf.Collections;

namespace Phoenix.Game.Card
{
    public class DataCenter
    {
        public Cproto.CharInfo charInfo = new Cproto.CharInfo();
        public List<Cproto.CharCard> cards = new List<Cproto.CharCard>();
        
        public long UId;

        public string SceneServer;
        public ulong SceneId;
        public int SceneCfgId;
        // 场景内玩家的id
        public int SceneAvatarId;

        public FightStat.FightStat fightStat = new FightStat.FightStat();


        private static DataCenter _it = new DataCenter();
        public static DataCenter It
        {
            get
            {
                return _it;
            }
        }

        public void Start()
        {
            bindEvents(true);
        }

        private void bindEvents(bool bind)
        {
            var events = GlobalEvents.It.events;
            
            events.Bind(EventDefine.CharInfo, onCharInfo, bind);
            events.Bind(EventDefine.LoginSucc, onLoginSucc, bind);
            events.Bind(EventDefine.RefreshServerCards, onRefreshCards, bind);
        }        

        private void onCharInfo(params object[] args)
        {
            var e = args[0] as HEventCharInfo;
            this.charInfo = e.charInfo.Clone();            
        }

        private void onRefreshCards(params object[] args)
        {
            var e = args[0] as HEventServerRefreshCards;
            copyCards(e.msg.Cards);
            // 刷界面
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventRefreshCards());
        }

        private void onLoginSucc(params object[] args)
        {
            var e = args[0] as HEventLoginSucc;
            UId = e.UId;
        }

        private void copyCards(pbc::RepeatedField<global::Cproto.CharCard> from) 
        {
            cards.Clear();
            for(var i = 0; i < from.Count; i ++)
            {
                cards.Add(from[i].Clone());
            }
        }

        public Cproto.CharCard GetCard(int id)
        {
            foreach(var one in cards)
            {
                if (one.Id == id)
                    return one;
            }
            return null;
        }
    }
}// namespace Phoenix
