using System.Text;

namespace Phoenix.Utils
{
    public static class DebugUtil
    {
        public static string DumpBuf(byte[] buf, int offset, int size, int maxDumpSize)
        {
            int useSize = size > maxDumpSize ? maxDumpSize : size;

            StringBuilder sb = new StringBuilder();
            for (var i = 0; i < useSize; i++)
            {
                sb.AppendFormat("{0:X} ", buf[i + offset]);
            }
            return sb.ToString();
        }
    }
} // Phoenix.Utils