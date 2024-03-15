using System;
using System.IO;
using System.Reflection;
using System.Text;
using UnityEditor;
using UnityEngine;
using UnityEngine.UI;

namespace Editor
{
    public class PanelCodeGen : EditorWindow
    {
        // [MenuItem("Tools/Panel管理")]
        // private static void ShowWindow()
        // {
        //     var window = GetWindow<PanelCodeGen>();
        //     window.titleContent = new GUIContent("Panel管理");
        //     window.Show();
        // }

        [MenuItem("Assets/生成Panel代码", false, 0)]
        private static void GenCode()
        {
            GameObject prefab = Selection.activeGameObject;

            if (!prefab.name.StartsWith("Panel"))
            {
                return;
            }

            if (FindScript(prefab.name) != null)
            {
                EditorUtility.DisplayDialog("提示", "代码已存在", "OK");
                return;
            }
            EditorUtility.DisplayProgressBar("提示","代码生成中……",0.6f);
            CodeGen(prefab);
            AssetDatabase.Refresh();
            EditorUtility.ClearProgressBar();
            EditorUtility.DisplayDialog("提示", "代码生成成功", "OK");
        }

        private static Type FindScript(string scriptName)
        {
            AppDomain appDomain = AppDomain.CurrentDomain;
            Assembly[] assemblies = appDomain.GetAssemblies();
            foreach (Assembly assembly in assemblies)
            {
                if (assembly.IsDynamic)
                {
                    continue;
                }

                Type[] types = assembly.GetExportedTypes();
                foreach (Type type in types)
                {
                    if (type.Name.Equals(scriptName))
                    {
                        return type;
                    }
                }
            }

            return null;
        }

        private static void CodeGen(GameObject prefab)
        {
            StreamWriter streamWriter = File.CreateText($"Assets/Scripts/UI/Panels/{prefab.name}.cs");
            streamWriter.AutoFlush = true;
            streamWriter.Write(GetSourceCodeFromTpl(prefab));
            streamWriter.Close();
            Debug.Log($"代码生成:{prefab.name}");
        }

        private static string GetRoute(string root, Transform transform, string splitter = "/")
        {
            string result = transform.name;
            Transform parent = transform.parent;
            while (!parent.name.Equals(root))
            {
                result = $"{parent.name}{splitter}{result}";
                parent = parent.parent;
            }

            return result;
        }

        private static string GetSourceCodeFromTpl(GameObject prefab)
        {
            string code = tpl;
            Button[] buttons = prefab.GetComponentsInChildren<Button>();
            Text[] texts = prefab.GetComponentsInChildren<Text>();

            if (buttons.Length == 0 && texts.Length == 0)
            {
                code = code.Replace("{Fields}", "\r\n");
                code = code.Replace("{FieldsFind}", "\r\n");
            }
            else
            {
                {
                    StringBuilder sb = new StringBuilder("\r\n");
                    foreach (Button button in buttons)
                    {
                        sb.AppendLine($@"        private Button _{button.name};");
                    }

                    foreach (Text text in texts)
                    {
                        if (text.name.Equals("Text"))
                        {
                            continue;
                        }

                        sb.AppendLine($@"        private Text _{text.name};");
                    }

                    code = code.Replace("{Fields}", sb.ToString());
                }

                {
                    //_btnSelectEnemy = TransformUtil.FindComponent<Button>(_root, "BG/btnSelectEnemy");
                    StringBuilder sb = new StringBuilder("\r\n");
                    foreach (Button button in buttons)
                    {
                        sb.AppendLine(
                            $@"           _{button.name} = TransformUtil.FindComponent<Button>(_root, ""{GetRoute(prefab.name, button.transform)}"");");
                    }

                    foreach (Text text in texts)
                    {
                        if (text.name.Equals("Text"))
                        {
                            continue;
                        }

                        sb.AppendLine(
                            $@"           _{text.name} = TransformUtil.FindComponent<Text>(_root, ""{GetRoute(prefab.name, text.transform)}"");");
                    }

                    code = code.Replace("{FieldsFind}", sb.ToString());
                }
            }

            code = code.Replace("{PrefabName}", prefab.name);

            return code;
        }

        private const string tpl = @"
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    [StringType(""{PrefabName}"")]
    public class {PrefabName} : BasePanel
    {
{Fields}        
        public override void OnReady()
        {
            SetDepth(PanelDepth.AboveNormal + 10);
            base.OnReady();

            BindEvents(true);
{FieldsFind}
        }

        public override void OnDestroy()
        {
            BindEvents(false);
        }

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;
        }
    }
} // namespace Phoenix
";
    }
}