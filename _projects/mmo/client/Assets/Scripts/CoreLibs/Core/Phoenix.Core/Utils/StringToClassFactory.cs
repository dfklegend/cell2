using System.Collections.Generic;

namespace Phoenix.Utils
{
    public class StringToClassFactory<TClass>
        where TClass:class
    {
        Dictionary<string, System.Type> _types = new Dictionary<string,System.Type>();

        public StringToClassFactory()
        {
            RegisterAll();
        }
       
        void AddType(System.Type type)
        {
            var key = StringType.GetValue(type);
            if (string.IsNullOrEmpty(key))
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

        public TClass Create(string script)
        {
            System.Type type;
            if (_types.TryGetValue(script, out type))
                return System.Activator.CreateInstance(type) as TClass;
            return null;
        }
    }
}