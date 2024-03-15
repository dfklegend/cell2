using System;

namespace Phoenix.Utils
{

    public static class ReflectUtil
    {
        // 
        public static bool HasOverrideFunc(Type type, Type excludeBaseType, string name)
        {
            var method = type.GetMethod(name);
            if (method == null)
                return false;
            // 有中间基类定义即可
            return method.IsVirtual && method.DeclaringType != excludeBaseType;
        }

        public static bool HasOverrideFunc(Type type, string name)
        {
            var method = type.GetMethod(name);
            if (method == null)
                return false;
            // 定义类型和当前类型一致
            return method.IsVirtual && method.DeclaringType == type;
        }
    }
} // Phoenix.Utils