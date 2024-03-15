using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;

namespace Phoenix.Game.Card
{	
    public class DisplayWorld : Singleton<DisplayWorld>
    {
        public float LogicToDisplayFactorW = Consts.GridSizeWidth;
        public float LogicToDisplayFactorH = Consts.GridSizeHeight;
        public const float ScreenWidth = 720f;
        public const float ScreenHeight = 1280f;
        // 2D世界视图
        public float windowCenterX = 0;
        public float windowCenterY = 0;
        int _focusId = -1;

        Grids grids = new Grids();

        public void SetFocusEntity(int id)
        {
            _focusId = id;
        }

        public Vector2 LogicToDisplay(Vector3 pos)
        {
            return new Vector2(pos.x * LogicToDisplayFactorW,
                pos.z * LogicToDisplayFactorH);
        }

        public Vector2 LogicToRelativeDisplay(Vector3 pos)
        {
            return new Vector2(pos.x * LogicToDisplayFactorW - windowCenterX,
                pos.z * LogicToDisplayFactorH - windowCenterY);
        }

        public void Init()
        {
            _focusId = -1;
            windowCenterX = 0;
            windowCenterY = 0;
            grids.Init(UIMgr.It.GetPanel<PanelFightChars>().GetGridRoot());
        }

        public void Destroy()
        {
            grids.Destroy();
            _focusId = -1;
            windowCenterX = 0;
            windowCenterY = 0;
        }

        public void Update()
        {            
            updateWindowByFocus();
            grids.Update(windowCenterX, windowCenterY);
            UnitModelMgr.It.Update();
        }

        private void updateWindowByFocus()
        {   
            Character focusChar = FightCtrl.It.GetChar(_focusId);
            if (_focusId == -1 || focusChar == null)
            {
                windowCenterX = 0f;
                windowCenterY = 0f;
                return;
            }

            Vector2 centerPos = LogicToDisplay(focusChar.pos);

            windowCenterX = centerPos.x;
            windowCenterY = centerPos.y;
        }
    }
} // namespace Phoenix
