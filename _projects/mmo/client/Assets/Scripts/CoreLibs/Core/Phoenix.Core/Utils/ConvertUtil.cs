namespace Phoenix.Utils
{   
    public class ConvertUtil
    {
        public static int ToInt(string inValue, int def = 0)
        {
            int value = def;
            if (int.TryParse(inValue, out value))
                return value;
            return def;
        }

        public static float ToFloat(string inValue, float def = 0f)
        {
            float value = def;
            if( float.TryParse(inValue, out value))
                return value;
            return def;
        }       
    }
}
