using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.EventSystems;

namespace Phoenix.Game
{
    // 用于保存style对象来处理点击事件
    // 处理条目的被点击事件
    public class ItemStyleHolder : MonoBehaviour, IPointerClickHandler
    {
        public BaseItemListStyle style;
        public IShowItem item;

        // Start is called before the first frame update
        void Start()
        {

        }

        // Update is called once per frame
        void Update()
        {

        }

        public void OnPointerClick(PointerEventData pointerEventData)
        {            
            Debug.Log(name + " Game Object Clicked!");

            if (style == null)
                return;
            style.OnClick(item);
        }
    }

}