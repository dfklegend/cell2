using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{        
    public partial class Entity
    {
        private RunEnv _env;
        public RunEnv env { get { return _env; } }
        public void SetEnv(RunEnv env)
        {
            _env = env;
        }
    }
}
