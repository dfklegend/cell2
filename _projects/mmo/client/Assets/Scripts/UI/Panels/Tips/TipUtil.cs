using Phoenix.Game.FightEmulator.BagSystem;
using Phoenix.Utils;
using System;
using System.Text;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    public static class TipUtil
    {
        public static PanelEquipTip ShowEquip(eEquipTipShowType showType,
            IShowItem item)
        {
            var panel = UIMgr.It.GetPanel<PanelEquipTip>();

            EquipTipData data = new EquipTipData();            
            data.showType = showType;
            data.item = item;
            panel.ShowTip(data);

            panel.ShowBtn(1);
            panel.SetBtnText(0, "关闭");
            panel.ResetHandlers();
            return panel;
        }

        public static PanelEquipTip ShowReplaceEquip(eEquipTipShowType showType,
            int slot, IShowItem item, Action replaceCB)
        {
            var panel = UIMgr.It.GetPanel<PanelEquipTip>();
            EquipTipData data = new EquipTipData();
            data.slot = slot;
            data.showType = showType;
            data.item = item;
            panel.ShowTip(data);

            panel.ShowBtn(2);
            panel.SetBtnText(0, "更换");
            panel.SetBtnText(1, "关闭");
            panel.ResetHandlers();
            panel.SetHandler(0, replaceCB);
            return panel;
        }

        public static PanelSkillTip ShowSkill(eSkillTipShowType showType,
            IShowItem item)
        {
            var panel = UIMgr.It.GetPanel<PanelSkillTip>();

            var data = new SkillTipData();
            data.showType = showType;
            data.item = item;
            panel.ShowTip(data);

            panel.ShowBtn(1);
            panel.SetBtnText(0, "关闭");
            panel.ResetHandlers();
            return panel;
        }

        public static PanelSkillTip ShowReplaceSkill(eSkillTipShowType showType,
            int slot, IShowItem item, Action replaceCB)
        {
            var panel = UIMgr.It.GetPanel<PanelSkillTip>();
            var data = new SkillTipData();
            data.slot = slot;
            data.showType = showType;
            data.item = item;
            panel.ShowTip(data);

            panel.ShowBtn(2);
            panel.SetBtnText(0, "更换");
            panel.SetBtnText(1, "关闭");
            panel.SetHandler(0, replaceCB);
            return panel;
        }
    }
} // namespace Phoenix
