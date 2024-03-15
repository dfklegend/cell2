using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;

namespace Phoenix.Game.Card
{	
    public class Grid
    {
        private GameObject _go;
        public GameObject go { get { return _go; } }
        private RectTransform _rect;

        public void Init(GameObject go)
        {
            _go = go;
            _rect = go.GetComponent<RectTransform>();
        }

        public void SetPos(float x, float y)
        {
            _rect.anchoredPosition = new Vector2(x, y);
        }
    }

    public class Grids
    {
        private Transform _root;
        GameObject _prefabPoint;

        const float GridSizeWidth = Consts.GridSizeWidth;
        const float GridSizeHeight = Consts.GridSizeHeight;
        const int GridsWidth = 20;
        const int GridsHeight = 20;

        Grid[] _grids;

        public void Init(Transform root)
        {
            _root = root;
            _prefabPoint = Resources.Load<GameObject>("Panels/prefabs/point");

            _grids = new Grid[GridsWidth * GridsHeight];
            initGrids();
            refreshGrids(0, 0);
        }

        public void Destroy()
        {
            int index = 0;
            for (var row = 0; row < GridsHeight; row++)
            {
                for (var col = 0; col < GridsWidth; col++, index++)
                {
                    GameObject.Destroy(_grids[index].go);
                }
            }
        }

        void initGrids()
        {
            int index = 0;
            for (var row = 0; row < GridsHeight; row++)
            {  
                for (var col = 0; col < GridsWidth; col++, index++)
                {
                    var go = createGo();
                    go.transform.SetParent(_root, false);

                    _grids[index] = new Grid();
                    _grids[index].Init(go);
                }                
            }
        }

        private GameObject createGo()
        {
            return GameObject.Instantiate(_prefabPoint);
        }

        void refreshGrids(float centerX, float centerY)
        {
            // 找合适的第一个位置
            float gridCenterX = (float)Math.Floor(centerX / GridSizeWidth) * GridSizeWidth - centerX;
            float gridCenterY = (float)Math.Floor(centerY / GridSizeHeight) * GridSizeHeight - centerY;
            float startX = gridCenterX - (GridsWidth / 2) * GridSizeWidth;
            float startY = gridCenterY - (GridsHeight / 2) * GridSizeHeight;

            float x;
            float y = startY;
            int index = 0;
            for (var row = 0; row < GridsHeight; row ++)
            {
                x = startX;
                for(var col = 0; col < GridsWidth; col ++, x += GridSizeWidth, index ++)
                {
                    _grids[index].SetPos(x, y);
                }
                y += GridSizeHeight;
            }
        }

        public void Update(float centerX, float centerY)
        {
            refreshGrids(centerX, centerY);
        }
    }
} // namespace Phoenix
