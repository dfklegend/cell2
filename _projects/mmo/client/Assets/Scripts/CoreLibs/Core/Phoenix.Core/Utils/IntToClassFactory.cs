using System.Collections.Generic;

namespace Phoenix.Utils
{
    public class IntToClassFactory<TClass>
        where TClass:class
    {
        Dictionary<int, System.Type> _types = new Dictionary<int,System.Type>();

        public IntToClassFactory()
        {
            RegisterAll();
        }
       
        void AddType(System.Type type)
        {
            var key = IntType.GetValue(type);
            if (key == -1)
                return;
            _types[key] = type;
        }

        public void RegisterAll()
        {
            List<System.Type> types = SystemUtil.GetAllClass<TClass>();
            if (types == null)
                return;
            for (int i = 0; i < types.Count; i++)
            {
                AddType(types[i]);
            }
        }

        public TClass Create(int define)
        {
            System.Type type;
            if (_types.TryGetValue(define, out type))
                return System.Activator.CreateInstance(type) as TClass;
            return null;
        }
    }
}