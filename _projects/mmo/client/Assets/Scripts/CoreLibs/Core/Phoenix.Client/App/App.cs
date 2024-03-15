using Phoenix.Core;

namespace Phoenix.Client
{
    public class BaseClientApp : Core.BaseApp
    {
        public override void OnPrepare()
        {
            AppInit.SetConsole(new UnityConsole());
            Log.LogEnv.SetRootPath(Utils.UDataDirectoryMgr.WTFGetFullNameInExternPath(""));
            AppInit.SetAppDir(Utils.UDataDirectoryMgr.GetRootPersistentPath());

            AppComponentMgr.Register<OneEnvAppComponent>();
        }

        public override void OnStart()
        {
        }

        public override void OnUpdate()
        {
        }

        public override void OnStop()
        {   
        }
    }
}
