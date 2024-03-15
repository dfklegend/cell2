using UnityEngine;

namespace Phoenix.Entity
{	
    public interface IEntity
    {
        void SetEntityID(int ID);
        int GetEntityID();
        
        /*
         * 创建过程
         *      CreateEntity
         *      AddComponent
         *      AddDataComponent
         *      OnInit()
         */
        void OnInit();
        void Destroy();
        void OnDestroy();        
       
        // 是否结束
        bool IsOver();

        T AddComponent<T>() where T : Component;
    }

    public interface IEntityWorld
    {
        IEntity CreateEntity();
        IEntity GetEntity(int ID);
        void DestroyEntity(int ID);        
        void Update();

        void Destroy();
        bool IsOver();
    }

    public interface IDataComponent
    {
        // 动态添加的，会在Entity 下次update时触发
        void OnInit();
        void OnDestroy();
    }

    public interface IComponent
    {

    }

    public interface IVisitor
    {
        void Visit(IEntity e);
    }
} // namespace Phoenix
