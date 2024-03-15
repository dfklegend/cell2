using UnityEngine;

namespace Phoenix.Utils
{
    public static class TransformUtil
    {
        public static T FindComponent<T>(Transform tran, string name)
            where T: class
        {
            var one = tran.Find(name);
            if (one == null)
                return default(T);
            return one.GetComponent<T>();
        }
    }
} // Phoenix.Utils