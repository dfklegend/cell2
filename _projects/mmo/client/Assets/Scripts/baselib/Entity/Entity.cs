using UnityEngine;

namespace Phoenix.Entity
{	
    public class Entity : MonoBehaviour, IEntity
    {
        private int _ID = -1;
        public int ID { get { return _ID; } }
        private int _worldID = -1;
        public int WorldID { get { return _worldID; } }
        private bool _over = false;

        // 对应一个逻辑对象
        // unit应该还是调整为component
        private ILogicUnit _unit;
        public ILogicUnit unit { get { return _unit; } }

        public void SetEntityID(int ID)
        {
            _ID = ID;
        }

        public int GetEntityID()
        {
            return _ID;
        }        

        public void SetWorldID(int worldID)
        {
            _worldID = worldID;
        }

        public void BindLogicUnit(ILogicUnit u) 
        { 
            _unit = u;
            u.SetEntity(this);
        }
        
        public void OnInit()
        {

        }

        public void Update()
        {
            _unit?.Update();
        }

        public void Destroy()
        {
            if (IsOver())
                return;
            setOver();
            _unit?.Destroy();
            GameObject.Destroy(gameObject);
        }
        
        public void OnDestroy()
        {

        }

        // 设置结束
        private void setOver()
        {
            _over = true;
        }
        // 是否结束
        public bool IsOver()
        {
            return _over;
        }

        public T AddComponent<T>() where T : Component
        {
            if (!gameObject)
                return default(T);
            return gameObject.AddComponent<T>();
        }        
    }
    
} // namespace Phoenix
