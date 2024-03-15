using System;

namespace Phoenix.Utils
{
    [AttributeUsage(AttributeTargets.Class, AllowMultiple = false, Inherited = false)]
    public class StringType : Attribute
    {
        public string define { get; set; }

        public StringType(string v)
        {
            this.define = v;
        }

        public static string GetValue( Type t)
        {
            if (!t.IsDefined(typeof(StringType), false))
                return string.Empty;
            StringType attribute = System.Attribute.GetCustomAttribute(t,
                typeof(StringType)) as StringType;
            if (attribute == null)
                return string.Empty;
            return attribute.define;
        }
    }   
}