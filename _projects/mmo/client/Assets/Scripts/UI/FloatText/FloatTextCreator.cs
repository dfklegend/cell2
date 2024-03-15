using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game.FText
{    
    public static class FloatTextCreator
    {
        public static FloatText CreateLineUp(Vector2 pos, string text, float life,
            float speed, Color color)
        {
            FloatText one = FloatTextMgr.It.Create();
            FloatTextShower shower = ShowerFactory.It.Create(ShowType.LinerUp, life, 0f);
            one.SetText(text);
            one.SetShower(shower);
            one.SetScreenPos(pos);
            one.SetColor(color);
            return one;
        }
    }
} // namespace Phoenix
