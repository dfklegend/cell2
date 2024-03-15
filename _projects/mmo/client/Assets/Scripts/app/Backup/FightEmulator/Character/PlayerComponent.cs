using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System.Text;


namespace Phoenix.Game.FightEmulator
{
    public class PlayerComponent : MonoBehaviour
    {
        private BagSystem.PlayerBags _bags;
        public BagSystem.PlayerBags bags { get { return _bags; } }
        public void SetBags(BagSystem.PlayerBags bags)
        {
            _bags = bags;
        }
    }

}// namespace Phoenix
