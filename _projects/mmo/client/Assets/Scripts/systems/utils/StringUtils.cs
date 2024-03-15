
namespace Phoenix.Game
{
    public static class StringsUtils
    {
        public static int GetIndexFromArray(string[] strings, string v)
        {
            for (var i = 0; i < strings.Length; i ++)
            {
                if (v == strings[i])
                    return i;
            }
            return -1;
        }
        
        public static string GetStringAtIndex(string[] strings, int index)
        {
            if (index < 0 || index >= strings.Length)
                return "";
            return strings[index];
        }
    }
    
}