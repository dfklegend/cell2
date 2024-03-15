using System.Collections;
using System.Collections.Generic;
using UnityEngine;

// 从动画节点copy位置变化
public class CopyRectPos : MonoBehaviour
{
    public Camera camera;
    public RectTransform src;
    public RectTransform tar;
    void Start()
    {
        
    }
 
    void LateUpdate()
    {
        if (src == null || tar == null)
            return;
        Vector3 pos = src.transform.position;
        RectTransform parentRect = tar.parent.GetComponent<RectTransform>();
        var newPos = Phoenix.Game.UIUtil.WorldToLocal(camera, pos, parentRect);
        tar.anchoredPosition = newPos;
    }
}
