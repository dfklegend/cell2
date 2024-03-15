using UnityEngine;
using System.Text;

namespace Phoenix.Utils
{
    // 负责可以存储文件的目录支持相关
    // 在不同环境文件路径支持
    public class UDataDirectoryMgr
    {

        // 根目录
        string _rootPersistentPath;
        // 根目录/UData
        string _persistentDataPath;
        string _streamingAssetsWWWRoot;
        private static UDataDirectoryMgr _instance;
        public static UDataDirectoryMgr It
        {
            get
            {
                if (null == _instance)
                {
                    new UDataDirectoryMgr();
                }
                return _instance;
            }
        }

        UDataDirectoryMgr()
        {
            _instance = this;

            InitAllPath();
        }

        #region Path
        static string initGetRootPath()
        {
            string path;
            switch (Application.platform)
            {
                case RuntimePlatform.WindowsEditor:
                case RuntimePlatform.OSXEditor:
                    path = Application.dataPath + "/../";
                    break;
                case RuntimePlatform.WindowsPlayer:
                case RuntimePlatform.OSXPlayer:
                    path = Application.dataPath + "/../";
                    break;
                case RuntimePlatform.IPhonePlayer:
                case RuntimePlatform.Android:
                    path = Application.persistentDataPath;
                    break;
                default:
                    path = Application.persistentDataPath;
                    break;
            }
            return SureCompletePath(path);
        }

        static string initGetDataPath()
        {
            return initGetRootPath() + "UData/";
        }

        public static string GetRootPersistentPath()
        {
            return It._getRootPersistentDataPath();
        }

        private string _getRootPersistentDataPath()
        {
            return _rootPersistentPath;
        }

        public static string GetPersistentDataPath()
        {
            return It._getPersistentDataPath();
        }

        private string _getPersistentDataPath()
        {
            return _persistentDataPath;
        }

        static string SureCompletePath(string path)
        {
            if (path.Length == 0)
                return string.Empty;
            if (path[path.Length - 1] != '/')
                path += "/";
            return path;
        }

        void InitAllPath()
        {
            // X
            _rootPersistentPath = initGetRootPath();
            // X/UData
            _persistentDataPath = initGetDataPath();
            PathUtil.SurePath(RootGetFullNameInExternPath("WTF/"));
            PathUtil.SurePath(RootGetFullNameInExternPath("UData/"));

            _streamingAssetsWWWRoot = Application.streamingAssetsPath;
            if (!_streamingAssetsWWWRoot.Contains("://"))
                _streamingAssetsWWWRoot = "file://" + _streamingAssetsWWWRoot;
            _streamingAssetsWWWRoot = SureCompletePath(_streamingAssetsWWWRoot);
        }


        // 直接文件系统
        // UData/...
        public static string GetFullNameInExternPath(string strFileName)
        {
            StringBuilder sb = new StringBuilder(255);
            sb.AppendFormat("{0}{1}", GetPersistentDataPath(), strFileName);
            return sb.ToString();
        }

        // ..
        public static string RootGetFullNameInExternPath(string strFileName)
        {
            StringBuilder sb = new StringBuilder(255);
            sb.AppendFormat("{0}{1}", GetRootPersistentPath(), strFileName);
            return sb.ToString();
        }

        // WTF/...
        public static string WTFGetFullNameInExternPath(string strFileName)
        {
            StringBuilder sb = new StringBuilder(255);
            sb.AppendFormat("{0}WTF/{1}", GetRootPersistentPath(), strFileName);
            return sb.ToString();
        }

        public static string GetWWWPathInStreamingAssets(string fn)
        {
            return string.Format("{0}{1}", It._streamingAssetsWWWRoot, fn);
        }
        #endregion

    }
}
