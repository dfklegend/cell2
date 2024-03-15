using System;

namespace Phoenix.Utils
{
    [AttributeUsage(AttributeTargets.Class, AllowMultiple = false, Inherited = false)]
    public class IntType : Attribute
    {
        public int define { get; set; }

        public IntType(int v)
        {
            this.define = v;
        }

        public static int GetValue( System.Type type)
        {
            if (!type.IsDefined(typeof(IntType), false))
                return -1;
            IntType attribute = System.Attribute.GetCustomAttribute(type,
                typeof(IntType)) as IntType;
            if (attribute == null)
                return -1;
            return attribute.define;
        }
    }   
}