using Phoenix.Core;
using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using UnityEngine;

namespace Phoenix.Res
{	
    public class ResourceMgr : Singleton<ResourceMgr>
    {
        public string LoadTextFile(string path)
        {
            var asset = Resources.Load<TextAsset>(path);
            if (asset != null)
                return asset.text;
            return "";
        }
    }
} // namespace Phoenix
