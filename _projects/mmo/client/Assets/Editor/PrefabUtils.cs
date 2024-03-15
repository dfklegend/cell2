using UnityEditor;
using UnityEditor.Experimental.SceneManagement;
using UnityEngine;

namespace Editor
{
    public static class PrefabUtils
    {
        [InitializeOnLoadMethod]
        private static void Load()
        {
            PrefabStage.prefabStageClosing -= OnPrefabClosed;
            PrefabStage.prefabStageClosing += OnPrefabClosed;

            PrefabStage.prefabStageOpened -= OnPrefabOpen;
            PrefabStage.prefabStageOpened += OnPrefabOpen;

        }

        private static void OnPrefabOpen(PrefabStage prefab)
        {
            // PanelCodeGen.CodeGen(prefab.prefabContentsRoot);
        }

        private static void OnPrefabClosed(PrefabStage prefab)
        {
            OnPrefabInstanceSaved(prefab);
        }

        private static void OnPrefabInstanceSaved(PrefabStage prefab)
        {
            if (!CheckValidPrefab(prefab.prefabContentsRoot))
            {
                EditorUtility.DisplayDialog("错误！", $"{prefab.prefabContentsRoot.name}检测到不合法预制体修改！", "好的");
            }
        }

        private static bool CheckValidPrefab(GameObject obj)
        {
            //TODO 添加prefab校验逻辑
            return true;
        }
    }
}