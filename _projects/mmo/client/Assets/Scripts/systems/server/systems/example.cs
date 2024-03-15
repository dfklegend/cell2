using System;
using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;


namespace Phoenix.Game
{
    public class ExampleSystem : BaseSystem
    {
        public override string GetName()
        {
            return "example";
        }

        public ExampleSystem()
        {
            registerCmd<Cproto.TestAddRet>("test1", onCmdTest1);
        }

        protected void onCmdTest1(object data)
        {
            Cproto.TestAddRet args = data as Cproto.TestAddRet;
            Log.LogCenter.Default.Debug($"onCmdTest1: {args.Result}");
        }
    }
}
