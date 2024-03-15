using System.Collections.Generic;

namespace Phoenix.Game
{
    public class IniFile
    {
        class IniSection
        {
            Dictionary<string, string> _kv = new Dictionary<string, string>();

            public string MakeContent()
            {
                string content = string.Empty;
                foreach (KeyValuePair<string, string> kv in _kv)
                {
                    content += SaveLine(kv.Key, kv.Value);
                }
                return content;
            }

            string SaveLine(string key, string value)
            {
                // 检查空格
                return string.Format("{0}={1}\r\n", key, value);
            }


            public string GetValue(string key, string def)
            {
                string value;
                if (_kv.TryGetValue(key, out value))
                    return value;
                return def;
            }

            public void SetValue(string key, string value)
            {
                _kv[key] = value;
            }

            public bool HasKey(string key)
            {
                return _kv.ContainsKey(key);
            }

            public bool GetBool(string key, bool def = false)
            {
                if (!HasKey(key))
                    return def;
                bool ret;
                if (bool.TryParse(_kv[key], out ret))
                    return ret;
                return def;
            }

            public void SetBool(string key, bool value)
            {
                SetValue(key, value.ToString());
            }

            public int GetInt(string key, int def = 0)
            {
                if (!HasKey(key))
                    return def;
                int ret;
                if (int.TryParse(_kv[key], out ret))
                    return ret;
                return def;
            }

            public void SetInt(string key, int value)
            {
                SetValue(key, value.ToString());
            }

            public float GetFloat(string key, float def = 0f)
            {
                if (!HasKey(key))
                    return def;
                float ret;
                if (float.TryParse(_kv[key], out ret))
                    return ret;
                return def;
            }

            public void SetFloat(string key, float value)
            {
                SetValue(key, value.ToString());
            }

            public string GetString(string key, string def = "")
            {
                return GetValue(key, def);
            }

            public void SetString(string key, string value)
            {
                SetValue(key, value);
            }           
        }

        Dictionary<string, IniSection> _sections = new Dictionary<string, IniSection>();

        public void Clear()
        {
            _sections.Clear();
        }

        public void Load(string file)
        {
            Clear();
            Utils.IniUtil.ParseIniFile(file, ConfigFileEntryParseCallback);
        }

        public void LoadFromContent(string content)
        {
            Clear();
            Utils.IniUtil.ParseConfigText(content, ConfigFileEntryParseCallback);
        }

        public void Save(string file)
        {
            string content = string.Empty;
            foreach (KeyValuePair<string, IniSection> kv in _sections)
            {
                if( !string.IsNullOrEmpty(kv.Key) )
                    content += string.Format("[{0}]\r\n", kv.Key);
                content += kv.Value.MakeContent();
            }
            Utils.PathUtil.SurePath(file);
            Res.ResUtil.WriteTextFile(file, content);
        }

        public void ConfigFileEntryParseCallback(string baseKey, string subKey, string val, object userData)
        {   
            GetSection(baseKey).SetValue(subKey, val);
        }

        IniSection GetSection(string baseKey, bool createIfMiss = true)
        {
            IniSection section;
            if (_sections.TryGetValue(baseKey, out section))
                return section;
            if (!createIfMiss)
                return null;
            section = new IniSection();
            _sections[baseKey] = section;
            return section;
        }

        public bool HasSection(string baseKey)
        {
            return GetSection(baseKey, false) != null;
        }

        public bool HasValue(string baseKey, string key)
        {
            IniSection section = GetSection(baseKey, false);
            if (section == null)
                return false;
            return section.HasKey(key);
        }

        public string GetValue(string baseKey, string key, string def)
        {
            IniSection section = GetSection(baseKey);
            if (section == null)
                return def;
            return section.GetValue(key, def);
        }

        public void SetValue(string baseKey, string key, string value)
        {
            IniSection section = GetSection(baseKey);
            section.SetValue(key, value);
        }

        public bool GetBool(string baseKey, string key, bool def = false)
        {
            if (!HasValue(baseKey, key))
                return def;
            return GetSection(baseKey).GetBool(key, def);
        }

        public void SetBool(string baseKey, string key, bool value)
        {
            SetValue(baseKey, key, value.ToString());
        }

        public int GetInt(string baseKey, string key, int def = 0)
        {
            if (!HasValue(baseKey, key))
                return def;
            return GetSection(baseKey).GetInt(key, def);
        }

        public void SetInt(string baseKey, string key, int value)
        {
            SetValue(baseKey, key, value.ToString());
        }

        public float GetFloat(string baseKey, string key, float def = 0f)
        {
            if (!HasValue(baseKey, key))
                return def;
            return GetSection(baseKey).GetFloat(key, def);
        }

        public void SetFloat(string baseKey, string key, float value)
        {
            SetValue(baseKey, key, value.ToString());
        }

        public string GetString(string baseKey, string key, string def = "")
        {
            return GetValue(baseKey, key, def);
        }

        public void SetString(string baseKey, string key, string value)
        {
            SetValue(baseKey, key, value);
        }
    }

    public class SIniFile
    {
        IniFile _ini = new IniFile();

        public void Load(string file)
        {
            _ini.Load(file);
        }

        public void LoadFromContent(string content)
        {
            _ini.LoadFromContent(content);
        }

        public void Clear()
        {
            _ini.Clear();
        }

        public void Save(string file)
        {
            _ini.Save(file);
        }
        
        public string GetValue(string key, string def)
        {
            return _ini.GetValue(string.Empty, key, def);
        }

        public void SetValue(string key, string value)
        {
            _ini.SetValue(string.Empty, key, value);
        }


        public bool HasKey(string key)
        {
            return _ini.HasValue(string.Empty, key);
        }

        public bool GetBool(string key, bool def = false)
        {
            return _ini.GetBool(string.Empty, key, def);
        }

        public void SetBool(string key, bool value)
        {
            SetValue(key, value.ToString());
        }

        public int GetInt(string key, int def = 0)
        {
            return _ini.GetInt(string.Empty, key, def);
        }

        public void SetInt(string key, int value)
        {
            SetValue(key, value.ToString());
        }

        public float GetFloat(string key, float def = 0f)
        {
            return _ini.GetFloat(string.Empty, key, def);
        }

        public void SetFloat(string key, float value)
        {
            SetValue(key, value.ToString());
        }

        public string GetString(string key, string def = "")
        {
            return GetValue(key, def);
        }

        public void SetString(string key, string value)
        {
            SetValue(key, value);
        }
    }
}
