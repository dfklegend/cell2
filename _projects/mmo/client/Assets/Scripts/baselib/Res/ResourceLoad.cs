using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using UnityEngine;


namespace Phoenix.Res
{	
    

    public class ResourceLoadType : IResLoadType
    {
        
        public void LoadPrefabAsync<T>(string path, Action<T> callback)
            where T: class
        {
            var req = Resources.LoadAsync(path, typeof(T));
            req.completed += (go) => 
            {
                callback(go as T);
            };
        }
    }
} // namespace Phoenix
