using System.Collections.Generic;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Log;
using System.Text;

namespace Phoenix.Game
{
    [StringType("PanelCardDetail")]
    public class PanelCardDetail : BasePanel
    {
        private Text _name;
        private Text _attrs;
        private List<CardEquipSlot> _equips = new List<CardEquipSlot>();
        private List<CardEquipSlot> _skills = new List<CardEquipSlot>();

        public override void OnReady()
        {
            SetDepth(PanelDepth.Normal + 100);
            base.OnReady();

            var btnClose = _root.Find("funcs/btnClose").GetComponent<Button>();
            btnClose.onClick.AddListener(() => {
                onClose();
            });

            _name = _root.Find("BG/fixed/base/name").GetComponent<Text>();
            _attrs = _root.Find("BG/scrollview/view/content/attrs/attrs").GetComponent<Text>();

            _name.text = "Î´Öª½ÇÉ«";
            _attrs.text = "ÊôÐÔÎ´Öª";
            initAllSlots();

            bindEvents(true);            
        }

        private void initAllSlots()
        {            
            initSlots(_root.Find("BG/scrollview/view/content/equips/items"), _equips, onClickEquip);
            initSlots(_root.Find("BG/scrollview/view/content/skills/items"), _skills, onClickSkill);
        }

        private void initSlots(Transform node, List<CardEquipSlot> slots, CardEquipSlot.OnClickSlot handler)
        {
            slots.Clear();
            for (var i = 0; i < CardConsts.MaxEquip; i++)
            {
                var one = new CardEquipSlot();
                one.Init(node.Find($"item{i}"), i, handler);
                slots.Add(one);
            }
        }

        public override void OnDestroy()
        {
            bindEvents(false);
        }

        private void bindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;            
        }

        private void showAttrs()
        {
            var cardSystem = Systems.It.GetSystem<CharCardSystem>("charcard");
            var sb = new StringBuilder();            
            addAttr(sb, cardSystem, eAttrType.HPMax);
            addAttr(sb, cardSystem, eAttrType.AttackSpeed);
            addAttr(sb, cardSystem, eAttrType.PhysicPower);
            addAttr(sb, cardSystem, eAttrType.MagicPower);
            addAttr(sb, cardSystem, eAttrType.PhysicArmor);
            addAttr(sb, cardSystem, eAttrType.MagicArmor);
            _attrs.text = sb.ToString();
        }

        private void addAttr(StringBuilder sb, CharCardSystem cardSystem, eAttrType index)
        {
            float value = 0;
            var one = cardSystem.GetAttr((int)index);
            if (one != null)
            {
                value = one.value;
            }            
            sb.AppendFormat("{0}: {1}\r\n", index.ToString(), value);
        }

        private void showBaseInfo()
        {
            var card = Systems.It.GetSystem<CharCardSystem>("charcard").GetCardDetail();
            if (card == null)
                return;
            _name.text = card.Name;
        }

        protected override void onShow()
        {
            base.onShow();
            refreshAll();
        }

        private void refreshAll()
        {
            showBaseInfo();
            showAttrs();
            showEquips();
            showSkills();
        }

        private void showEquips()
        {
            var card = Systems.It.GetSystem<CharCardSystem>("charcard").GetCardDetail();
            if (card == null)
                return;
            if (card.Equips == null)
                return;
            for(var i = 0; i < CardConsts.MaxEquip; i ++)
            {
                var slot = _equips[i];
                slot.RefreshInfo(null);
                if (i >= card.Equips.Count)
                    continue;
                var item = card.Equips[i];
                if (string.IsNullOrEmpty(item.EquipId))
                {
                    continue;
                }
                    
                slot.RefreshInfo(new Card.EquipItem(item.EquipId));
            }
        }

        private void showSkills()
        {
            var card = Systems.It.GetSystem<CharCardSystem>("charcard").GetCardDetail();
            if (card == null)
                return;
            if (card.Skills == null)
                return;
            for (var i = 0; i < CardConsts.MaxSkill; i++)
            {
                var slot = _skills[i];
                slot.RefreshInfo(null);
                if (i >= card.Skills.Count)
                    continue;
                var item = card.Skills[i];
                if (string.IsNullOrEmpty(item.SkillId))
                {
                    continue;
                }

                slot.RefreshInfo(new Card.SkillItem(item.SkillId));
            }
        }

        private string getCardEquip(int index)
        {
            var card = Systems.It.GetSystem<CharCardSystem>("charcard").GetCardDetail();
            if (card == null)
                return null;
            if (card.Equips == null)
                return null;
            if (index < 0 || index >= card.Equips.Count)
                return null;
            var item = card.Equips[index];
            if (string.IsNullOrEmpty(item.EquipId))
            {
                return null;
            }
            return item.EquipId;
        }

        private void onClickEquip(CardEquipSlot slot)
        {
            LogCenter.Default.Debug("click equip {0}", slot.GetIndex());
            var index = slot.GetIndex();
            string equipId = getCardEquip(index);
            if(equipId != null)
            {
                TipUtil.ShowReplaceEquip(eEquipTipShowType.Equip, index, new Card.EquipItem(equipId),
                    () => {
                        selectToReplace(index);
                    });
            } 
            else
            {
                // select to equip
                BagCtrl.It.SelectEquip(index, (result) =>
                {
                    BagCtrl.It.HideSelectPanel();
                    if (result.selected != null)
                    {
                        var select = result.selected.GetItemId();
                        Debug.Log($"Select: {select}");
                        // equip
                        Systems.It.GetSystem<CharCardSystem>("charcard").SetEquip(index, select, (code) => 
                        {
                            // refresh this 
                            refreshAll();
                        });
                    }
                });
            }
        }

        private void selectToReplace(int index)
        {
            BagCtrl.It.SelectEquipToReplace(index, null, (result) =>
            {
                BagCtrl.It.HideSelectPanel();
                string selected = "";
                if (result.selected != null)
                {
                    selected = result.selected.GetItemId();
                }
                if(result.selectUnequip)
                {
                    // Ð¶ÔØ
                    selected = "";
                }

                Debug.Log($"Select: {selected}");
                // equip
                Systems.It.GetSystem<CharCardSystem>("charcard").SetEquip(index, selected, (code) =>
                {
                    // refresh this 
                    refreshAll();
                });
            });
        }

        private void onClickSkill(CardEquipSlot slot) 
        {
            LogCenter.Default.Debug("click skill {0}", slot.GetIndex());

            var index = slot.GetIndex();
            string skillId = getCardSkill(index);
            if (skillId != null)
            {
                var curSkill = new Card.SkillItem(skillId);
                TipUtil.ShowReplaceSkill(eSkillTipShowType.Equip, index, curSkill,
                    () =>
                    {
                        selectToReplaceSkill(index, curSkill);
                    });
            }
            else
            {
                // select to equip
                BagCtrl.It.SelectSkill(index, (result) =>
                {
                    BagCtrl.It.HideSelectPanel();
                    if (result.selected != null)
                    {
                        var selected = result.selected.GetItemId();
                        Debug.Log($"Select: {selected}");

                        Systems.It.GetSystem<CharCardSystem>("charcard").SetSkill(index, selected, (code) =>
                        {
                            // refresh this 
                            refreshAll();
                        });
                    }
                });
            }
        }

        private void selectToReplaceSkill(int index, IShowItem item)
        {
            BagCtrl.It.SelectSkillToReplace(index, item, (result) =>
            {
                BagCtrl.It.HideSelectPanel();
                string selected = "";
                if (result.selected != null)
                {
                    selected = result.selected.GetItemId();
                }
                if (result.selectUnequip)
                {
                    // Ð¶ÔØ
                    selected = "";
                }

                Debug.Log($"Select: {selected}");
                // equip
                Systems.It.GetSystem<CharCardSystem>("charcard").SetSkill(index, selected, (code) =>
                {
                    // refresh this 
                    refreshAll();
                });
            });
        }

        private string getCardSkill(int index)
        {
            var card = Systems.It.GetSystem<CharCardSystem>("charcard").GetCardDetail();
            if (card == null)
                return null;
            if (card.Skills == null)
                return null;
            if (index < 0 || index >= card.Skills.Count)
                return null;
            var item = card.Skills[index];
            if (string.IsNullOrEmpty(item.SkillId))
            {
                return null;
            }
            return item.SkillId;
        }

        private void onClose()
        {
            Systems.It.GetSystem<CharCardSystem>("charcard").SaveCurCard();
            Hide();
        }
    }
} // namespace Phoenix
