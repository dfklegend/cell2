using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System.Text;

namespace Phoenix.Game.FightEmulator
{
    public partial class Character : IOwnerImpl, ICharacter
    {

        public string MakeBrief()
        {
            StringBuilder sb = new StringBuilder();
            
            sb.AppendLine($"名字: {name}");
            sb.AppendLine($"level: {level}");
            AttrsUtil.MakeAttrsBrief(sb, attrs);

            _equips.MakeBrief(sb);


            // 主手武器
            briefAddWeapon(sb, "主手", GetMainWeapon());
            briefAddWeapon(sb, "副手", GetOffHandWeapon());
            // 副手武器
            return sb.ToString();
        }

        private void briefAddWeapon(StringBuilder sb, string which, IWeapon w)
        {
            EquipUtil.briefAddWeapon(sb, which, w);
        }
    }

}// namespace Phoenix
