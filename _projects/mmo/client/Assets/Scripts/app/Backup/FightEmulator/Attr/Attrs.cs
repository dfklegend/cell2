using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{	
    public class Attrs
    {
        private Dictionary<string, Attr> _attrs = new Dictionary<string, Attr>();
        

        public Attrs()
        {            
        }

        public void NewAttr(string attrName)
        {
            _attrs[attrName] = new Attr();
        }

        public Attr GetAttr(string attrName)
        {
            Attr ret;
            if (!_attrs.TryGetValue(attrName, out ret))
                return null;
            return ret;
        }
        
        public void Reset()
        {
            var names = _attrs.Keys;
            foreach( var key in names)
            {
                _attrs[key].Reset();
            }
        }

        public void Dump()
        {
            Log.LogCenter.Default.Debug("---- attrs:");
            var names = _attrs.Keys;
            foreach (var key in names)
            {
                Log.LogCenter.Default.Debug($"{key}: {_attrs[key].final}");
            }
            Log.LogCenter.Default.Debug("---- end");
        }
    }
}// namespace Phoenix
