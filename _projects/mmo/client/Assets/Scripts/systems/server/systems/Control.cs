using System;
using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;


namespace Phoenix.Game
{
    public class ControlSystem : BaseSystem
    {
        public override string GetName()
        {
            return "control";
        }

        public ControlSystem()
        {
            registerCmd<Cproto.TestAddRet>("init", onCmdInit);
            registerCmd<Cproto.TestAddRet>("systeminfo", onCmdInfo);
        }

        protected void onCmdTest1(object data)
        {
            Cproto.TestAddRet args = data as Cproto.TestAddRet;
            Log.LogCenter.Default.Debug($"baseinfo.onCmdTest1: {args.Result}");
        }

        protected void onCmdInit(object data)
        {
            var args = data as Cproto.DaySign;
            //Log.LogCenter.Default.Debug($"baseinfo.onCmdInit: {args.Result}");
        }

        protected void onCmdInfo(object data)
        {
            Cproto.SystemInfo args = data as Cproto.SystemInfo;
            Log.LogCenter.Default.Debug($"baseinfo.onCmdInfo: {args.Type} - {args.Info}");
        }
    }
}
