using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using UnityEngine.UI;
using Phoenix.Utils;

namespace Phoenix.Game.Card
{	
    public class UnitModel : BaseModel
    {
        const int HP_WIDTH = 100;
        private GameObject _go;
        private Transform _root;
        private RectTransform _rect;

        private Image _icon;
        private Text _name;
        private Image _hp;
        private Image _mp;

        private Transform _animRoot;
        private Animator _animator;
        private AnimCtrl _animCtrl = new AnimCtrl();
        private Character _owner;
        

        public void Init(GameObject go, CharacterUnit unit)
        {
            if (!go)
                return;
            _go = go;
            _root = go.transform;
            _rect = go.GetComponent<RectTransform>();

            initCtrls();
            setInfo(unit);            
        }

        public override Transform GetTransform()
        {
            return _root;
        }

        private void initCtrls()
        {
            _icon = TransformUtil.FindComponent<Image>(_root, "visual/icon");
            _name = TransformUtil.FindComponent<Text>(_root, "visual/name/text");
            _hp = TransformUtil.FindComponent<Image>(_root, "visual/hp/value");
            _mp = TransformUtil.FindComponent<Image>(_root, "visual/mp/value");
            _animRoot = _root.Find("anim");
            _animator = TransformUtil.FindComponent<Animator>(_root, "anim/model");
            _animCtrl.Init(_animator);
        }

        private void setInfo(CharacterUnit unit)
        {
            _owner = unit.character;
            _name.text = unit.character.name;

            UpdateMP(unit.character.GetMPPercent());
            UpdateHP(unit.character.GetHPPercent());
            UpdatePos();
        }

        public void InitInfo(Character c)
        {
            _name.text = c.name;
        }

        public override void SetPos(float x, float y)
        {
            _rect.anchoredPosition = new Vector2(x, y);           
        }

        public Vector2 GetPos()
        {
            return _rect.anchoredPosition;
        }

        public Vector3 GetWorldPos()
        {
            return _rect.position;
        }

        public void UpdateHP(float percent)
        {
            const float HEI = 20;
            _hp.rectTransform.sizeDelta = new Vector2(HP_WIDTH * percent, HEI);
        }

        public void UpdateMP(float percent)
        {
            const float HEI = 10;
            _mp.rectTransform.sizeDelta = new Vector2(HP_WIDTH * percent, HEI);
        }

        public void PlayAnim(string anim, float time)
        {
            _animCtrl.PlayAnim(anim, time);
        }

        public void StopCurAnim()
        {
            _animCtrl.StopCurAnim();
        }

        public override void Update()
        {
            _animCtrl.Update();
            UpdatePos();
        }

        public void SetCharPos(Vector3 pos)
        {
            var uiPos = DisplayWorld.It.LogicToRelativeDisplay(pos);
            SetPos(uiPos.x, uiPos.y);
            updateDepth();
        }

        public void UpdatePos()
        {
            // 根据位置刷新界面位置
            Vector3 pos = _owner.pos;
            SetCharPos(pos);
        }

        private int getFocusSide()
        {
            return FightCtrl.It.GetPlayerSide();
        }

        // 下面，右边的在上面
        // 玩家对象最优先
        private void updateDepth()
        {
            var pos = GetPos();
            // Y轴负向
            // 本方阵营
            var factor = -pos.y * 1f + pos.x / 1000;
            if (_owner.side == getFocusSide())
                factor += 1000;
            _depth = factor;
        }

        public void LookAtTar(int tarId)
        {
            var tar = UnitModelMgr.It.GetModel(tarId);
            if (tar == null)
                return;
            Vector3 src = GetWorldPos();
            Vector3 dst = tar.GetWorldPos();

            Vector3 dir = dst - src;
            dir.Normalize();
            if (dir.sqrMagnitude < 1f)
                return;
            var angel = -(180f/Math.PI)*Math.Atan2(dir.x, dir.y);
            _animRoot.localEulerAngles = new Vector3(0, 0, (float)angel);
        }

        public override void Destroy()
        {
            if (!_go)
                return;
            GameObject.Destroy(_go);
        }

        public void SetIcon(string path)
        {
            var sprite = Resources.Load<Sprite>(path);
            if (sprite == null)
                return;
            _icon.sprite = sprite;
        }
    }    
} // namespace Phoenix
