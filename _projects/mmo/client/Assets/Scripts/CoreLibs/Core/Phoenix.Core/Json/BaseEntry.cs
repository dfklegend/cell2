using System.Collections;
using Newtonsoft.Json.Linq;

namespace Phoenix.Core.Json
{
    // 依次组织成节点树
    public abstract class BaseEntry
    {
        public virtual void LoadFromJson(JToken node)
        {
            JsonUtil.LoadValues(this, node);
        }

        public virtual void SaveToJson(JToken node)
        {
            JsonUtil.SaveValues(this, node);
        }
    }
}
