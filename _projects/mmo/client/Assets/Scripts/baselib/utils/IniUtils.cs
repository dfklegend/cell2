using System.Collections.Generic;

namespace Phoenix.Game
{
    public static class IniUtils
    {
        // load from Resources/
        public static void LoadFromResource(IniFile cfg, string path)
        {
            cfg.LoadFromContent(Res.ResourceMgr.It.LoadTextFile(path));
        }

        // load from UData/...
        public static void LoadFromUData(IniFile cfg, string path)
        {
            cfg.Load(Utils.UDataDirectoryMgr.GetFullNameInExternPath(path));
        }
    }
}
