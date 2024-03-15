using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{
    public partial class Entity
    {
        private string _name = "";
        public string name { get { return _name; } }

        public void SetName(string name)
        {
            _name = name;
        }

        // 从当前节点查找子节点
        // xx/xx
        public Entity Find(string path)
        {
            string[] subs = path.Split('/');
            // 找到对应的child
            if (subs.Length == 0)
                return null;
            return find(0, subs);
        }

        private Entity find(int head, string[] subs)
        {            
            if (subs.Length <= head)
                return null;
            var childName = subs[head];
            var child = findChild(childName);
            if (child == null)
                return null;
            if (subs.Length == head+1)
                return child;
            return child.find(head+1, subs);
        }

        private Entity findChild(string name)
        {
            for (var i = 0; i < _childs.Count; i++)
            {
                if (name == _childs[i].name)
                    return _childs[i];
            }
            return null;
        }
    }
}
