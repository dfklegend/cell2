using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using Phoenix.Core;
using Phoenix.Res;
using UnityEngine;

namespace Phoenix.Platform
{
    // 可以定义一些平台有关的行为
    public interface IPlatform
    {
        void Init();
    }

    public class DefaultPlatform : IPlatform
    {
        public void Init() { }
    }

    public class MobilePlatform : IPlatform
    {
        public void Init() 
        {
            csv.CSVEnv.loader = new csv.ResourceLoader();
        }
    }

    public class PlatformMgr : Singleton<PlatformMgr>
    {
        IPlatform _platform;

        public IPlatform platform 
        {
            get { return _platform; }
        }

        public void SetPlatform(IPlatform platform) 
        {
            _platform = platform;
        }
    }

    public static class PlatformUtils
    {
        private static IPlatform createByPlatform()
        {
#if UNITY_ANDROID || UNITY_IPHONE
            return new MobilePlatform();
#endif
            return new DefaultPlatform();
        }

        public static void InitPlatform()
        {
            IPlatform platform = createByPlatform();

            PlatformMgr.It.SetPlatform(platform);
            platform.Init();
        }
    }
}