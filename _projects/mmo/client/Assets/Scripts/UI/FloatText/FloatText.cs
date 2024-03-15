using System.Collections.Generic;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game.FText
{
    public interface FloatTextShower
    {
        void Create(params object[] args);
        void Update(FloatText text);
        bool IsOver();
    }
    public class FloatText
    {
        Transform _root;
        RectTransform _rect;
        Text _text;
        FloatTextShower _shower;
        Animator _animator;

        public void Init(GameObject go)
        {
            if (!go)
                return;
            _root = go.transform;
            _rect = go.GetComponent<RectTransform>();

            _rect.anchoredPosition = Vector2.zero;
            _root.localScale = Vector3.one;

            _text = TransformUtil.FindComponent<Text>(_root, "Text");
            _animator = TransformUtil.FindComponent<Animator>(_root, "Text");
            
        }

        public void Destroy()
        {
            if (!_root)
                return;
            Object.Destroy(_root.gameObject);
        }

        public Vector2 GetPos()
        {
            return _rect.anchoredPosition;
        }

        public void SetPos(Vector2 pos)
        {
            _rect.anchoredPosition = pos;
        }

        public void OffsetPos(Vector2 pos)
        {
            _rect.anchoredPosition += pos;
        }

        public void SetScreenPos(Vector2 pos)
        {
            var pos1 = UIUtil.ScreenToLocal(null, pos, _rect);
            SetPos(pos1);
        }

        public void SetColor(Color c)
        {
            _text.color = c;
        }

        public void SetText(string text)
        {
            _text.text = text;
        }

        public void SetShower(FloatTextShower shower)
        {
            _shower = shower;
        }

        public void Update()
        {
            _shower.Update(this);
        }

        public bool IsOver()
        {
            return _shower.IsOver();
        }

        public void PlayAnim(string anim)
        {
            _animator.Play(anim);
        }
    }
} // namespace Phoenix
