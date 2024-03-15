using System;
using System.Threading.Tasks;
using UnityEngine;


namespace Phoenix.Res
{
	// 提供统一载入接口
    // Resources AB
    public interface IResLoadType
    {
        void LoadPrefabAsync<T>(string path, Action<T> callback) where T : class;
    }
} // namespace Phoenix
