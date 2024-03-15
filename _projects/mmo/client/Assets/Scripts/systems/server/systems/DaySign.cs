using System;
using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;


namespace Phoenix.Game
{
    public class DaySignSystem : BaseSystem
    {
        protected bool _signed = false;
        public bool signed => _signed;
        public override string GetName()
        {
            return "daysign";
        }

        public DaySignSystem()
        {
            registerCmd<Cproto.DaySign>("init", onCmdInit);            
        }        

        protected void onCmdInit(object data)
        {
            var args = data as Cproto.DaySign;
            Log.LogCenter.Default.Debug($"DaySign.onCmdInit: {args.Signed}");
            _signed = args.Signed;
        }       

        public void Sign(Action<bool> cb)
        {
            this.Request("sign", null, (code) =>
            {
                if(code == 0) 
                {
                    _signed = true;
                    cb?.Invoke(true);
                    return;
                }
                cb?.Invoke(false);
            });
        }
    }    
}
