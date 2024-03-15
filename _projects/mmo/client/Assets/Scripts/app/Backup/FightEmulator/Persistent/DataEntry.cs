using System;
using Newtonsoft.Json.Linq;
using Phoenix.Core.Json;
using Phoenix.Core;

namespace Phoenix.Game.FightEmulator
{
    public class EnemyEntry : BaseEntry
    {
        public string cfgId = "";
        public int level = 1;

        public override void LoadFromJson(JToken node) 
        {
            cfgId = JsonUtil.AsStr(node, "cfgId", "");
            level = JsonUtil.AsInt(node, "level", 1);
        }

        public override void SaveToJson(JToken node)
        {
            node["cfgId"] = cfgId;
            node["level"] = level;
        }
    }

    public class PlayerEntry: BaseEntry
    {
        public int level = 1;
        public string mainHand = "";
        public string offHand = "";

        public override void LoadFromJson(JToken node)
        {
            level = JsonUtil.AsInt(node, "level", 1);
            mainHand = JsonUtil.AsStr(node, "mainHand", "");
            offHand = JsonUtil.AsStr(node, "offHand", "");
        }

        public override void SaveToJson(JToken node)
        {            
            node["level"] = level;
            node["mainHand"] = mainHand;
            node["offHand"] = offHand;
        }        
    }

    public class UserPersistData : Singleton<UserPersistData>
    {
        private const string FileName = "userdata.txt";
        private NodeEntry _root = new NodeEntry();
        public UserPersistData()
        {
            _root.AddEntry("player", new PlayerEntry());
            _root.AddEntry("enemy", new EnemyEntry());
        }

        public void LoadFromFile()
        {
            var path = Utils.UDataDirectoryMgr.GetFullNameInExternPath(FileName);
            var str = Res.ResUtil.LoadTextFile(path);
            _root.FromStr(str);
        }

        public bool IsValid()
        {
            return !string.IsNullOrEmpty(GetEntry<EnemyEntry>("enemy").cfgId);
        }

        public void SaveToFile()
        {
            var path = Utils.UDataDirectoryMgr.GetFullNameInExternPath(FileName);
            var str = _root.ToStr();
            Res.ResUtil.WriteTextFile(path, str);
        }

        public T GetEntry<T>(string name)
            where T : BaseEntry
        {
            return _root.GetEntry<T>(name);
        }
    }
} // Phoenix.Utils