using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using Phoenix.Res;
using UnityEngine;

namespace Phoenix.csv
{
    public interface ILoader
    {
        string LoadContent(string path);
    }

    public class FileLoader : ILoader
    {
        public string LoadContent(string path)
        {
            path = path + ".csv";
            return ResUtil.LoadTextFile(Utils.UDataDirectoryMgr.GetFullNameInExternPath(path));
        }
    }

    public class ResourceLoader : ILoader
    {
        public string LoadContent(string path)
        {
            return ResourceMgr.It.LoadTextFile(path);
        }
    }


    public static class CSVEnv
    {
        // android 直接从resouces读取，避免需要更新数据
        public static ILoader loader = new FileLoader();
    }
}