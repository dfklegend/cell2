using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;

namespace Phoenix.Game
{
    public class Config : Singleton<Config>
    {
        private IniFile _base = new IniFile();
        public IniFile baseCfg
        {
            get
            {
                return _base;
            }
        }

        private IniFile _main = new IniFile();
        public IniFile main
        {
            get 
            {
                return _main;
            }
        }

        public void Init()
        {
            IniUtils.LoadFromResource(_base, "ASys/base");
            IniUtils.LoadFromUData(_main, "cfg.txt");            
        }

        
    }
}

