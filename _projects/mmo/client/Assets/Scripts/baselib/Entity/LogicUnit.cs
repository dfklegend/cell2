using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;

namespace Phoenix.Entity
{	
    // 基础的逻辑单元
    // 每个Entity对应一个
    public interface ILogicUnit
    {
        void Update();
        int GetUnitType();
        void SetEntity(IEntity e);
        void Destroy();
    }
} // namespace Phoenix
