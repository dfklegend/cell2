using System;
using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;
using Phoenix.Game.Card;


namespace Phoenix.Game
{
    public class BaseInfoSystem : BaseSystem
    {
        public override string GetName()
        {
            return "baseinfo";
        }

        public BaseInfoSystem()
        {
            //registerCmd<Cproto.DaySign>("init", onCmdInit);
            registerCmd<Cproto.SystemInfo>("info", onCmdInfo);
            registerCmd<Cproto.CharInfo>("updatecharinfo", onCmdUpdateCharInfo);
        }

        protected void onCmdTest1(object data)
        {
            Cproto.TestAddRet args = data as Cproto.TestAddRet;
            Log.LogCenter.Default.Debug($"baseinfo.onCmdTest1: {args.Result}");
        }

        protected void onCmdInit(object data)
        {
            //var args = data as Cproto.DaySign;
            //Log.LogCenter.Default.Debug($"baseinfo.onCmdInit: {args.Result}");
        }

        protected void onCmdInfo(object data)
        {
            Cproto.SystemInfo args = data as Cproto.SystemInfo;
            Log.LogCenter.Default.Debug($"baseinfo.onCmdInfo: {args.Type} - {args.Info}");
        }

        protected void onCmdUpdateCharInfo(object data)
        {            
            var args = data as Cproto.CharInfo;
            Log.LogCenter.Default.Debug($"baseinfo.onCmdUpdateCharInfo: {args}");
            HEventUtil.Dispatch(GlobalEvents.It.events, new HEventCharInfo(args));
        }
    }    
}
