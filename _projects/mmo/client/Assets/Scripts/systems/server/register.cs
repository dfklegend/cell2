using System;
using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;

namespace Phoenix.Game
{
    public static class SystemsUtils
    {
        public static void RegisterAll(Systems systems)
        {
            systems.AddSystem(new ExampleSystem());
            systems.AddSystem(new CharCardSystem());
            systems.AddSystem(new BaseInfoSystem());
            systems.AddSystem(new ControlSystem());
            systems.AddSystem(new DaySignSystem());
        }
    }
}
