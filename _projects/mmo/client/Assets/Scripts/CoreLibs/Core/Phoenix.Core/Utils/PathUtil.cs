using System;
using System.IO;
using System.Collections.Generic;

namespace Phoenix.Utils
{
    public class PathUtil
    {
        public static void SurePath(string path)
        {
            string strDir = Path.GetDirectoryName(path);
            if (Directory.Exists(strDir))
                return;
            // 会创建多级目录
            Directory.CreateDirectory(strDir);
        }

        public static List<string> GetFilesInDir(string path)
        {
            List<string> ret = new List<string>();
            DirectoryInfo dir = new DirectoryInfo(path);
            if (!dir.Exists)
                return ret;
            GetFilesInDir(dir, ret);
            return ret;
        }

        static void GetFilesInDir(DirectoryInfo dir, List<string> ret)
        {
            FileInfo[] infos = dir.GetFiles();
            foreach (FileInfo fileInfo in infos)
            {
                ret.Add(fileInfo.FullName);
            }

            DirectoryInfo[] dirs = dir.GetDirectories();
            foreach (DirectoryInfo one in dirs)
            {
                GetFilesInDir(one, ret);
            }
        }
    }
}