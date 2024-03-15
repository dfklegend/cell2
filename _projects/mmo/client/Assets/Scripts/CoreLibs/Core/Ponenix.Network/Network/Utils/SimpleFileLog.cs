using System.IO;

namespace Phoenix.Utils
{
    public static class SimpeFileLogUtil
    {
        public static void AppendTextFile(string path, string content)
        {
            StreamWriter sw = new StreamWriter(path, true, System.Text.Encoding.UTF8);
            sw.Write(content);
            sw.Close();
        }
    }
} // Phoenix.Utils