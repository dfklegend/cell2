using UnityEngine;
using System.Collections.Generic;
using System.Reflection;
using Phoenix.csv;

namespace Phoenix.Game.FightEmulator
{
    // string到int的对应
    public class StringToInt
    {
        private Dictionary<string, int> _map = new Dictionary<string, int>();
        private int _begin;

        public void InitFromStrings(string[] strs, int begin = 0)
        {
            _begin = begin;
            
            for(var i = 0; i < strs.Length; i ++)
            {
                _map[strs[i]] = begin + i;
            }                
        }

        public int ToInt(string str)
        {
            int ret;
            if (_map.TryGetValue(str, out ret))
                return ret;
            return _begin;
        }
    }    
}
