using System.Collections.Generic;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    public static class UIUtil
    {
        //
        public static Vector2 WorldToLocal(Camera camera, Vector3 worldPos, RectTransform rect)
        {
            if (camera == null)
                camera = UIMgr.It.camera;
            var screenPos = WorldToScreen(camera, worldPos);
            return ScreenToLocal(camera, screenPos, rect);
        }

        // 世界坐标转化为 屏幕坐标
        public static Vector2 WorldToScreen(Camera camera, Vector3 worldPos)
        {
            if (camera == null)
                camera = UIMgr.It.camera;
            return RectTransformUtility.WorldToScreenPoint(camera, worldPos);            
        }

        // 屏幕坐标转换为 本地坐标
        public static Vector2 ScreenToLocal(Camera camera, Vector2 screenPos, RectTransform rect)
        {
            if (camera == null)
                camera = UIMgr.It.camera;
            Vector2 localPos;
            RectTransformUtility.ScreenPointToLocalPointInRectangle(rect, screenPos, camera, out localPos);
            return localPos;
        }

        public static Sprite LoadIcon(string icon)
        {
            return LoadSprite($"icons/{icon}");
        }

        public static Sprite LoadItemIcon(string icon)
        {
            return LoadSprite($"icons/items/{icon}");            
        }

        public static Sprite LoadSprite(string path)
        {           
            var sprite = Resources.Load<Sprite>(path);
            if (sprite == null)
            {
                Debug.LogError($"load {path} failed!");
            }
            return sprite;
        }
    }
} // namespace Phoenix
