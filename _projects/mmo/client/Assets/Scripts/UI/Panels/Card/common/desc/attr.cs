using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Log;
using System.Text;

// Éú³ÉÃèÊöÎÄ×Ö
namespace Phoenix.Game.Card
{
    public static class AttrFormatter
    {
        public static void FormatAttr(StringBuilder sb, string prefix, eAttrType attr, bool percent, float v)
        {
            if (attr == eAttrType.Invalid)
                return;
            if(!percent)
                sb.AppendFormat("{0}{1}: {2}\n", prefix,
                    StringsUtils.GetStringAtIndex(Strings.ATTR_NAMES, (int)attr), v);
            else
                sb.AppendFormat("{0}{1}: {2}%\n", prefix,
                    StringsUtils.GetStringAtIndex(Strings.ATTR_NAMES, (int)attr), v*100);
        }
    }
} // namespace Phoenix
