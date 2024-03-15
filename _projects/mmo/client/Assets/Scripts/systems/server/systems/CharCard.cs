using System;
using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;


namespace Phoenix.Game
{
    public class AttrItem
    {
        public int index = 0;
        public float value = 0;
    }


    public class CharCardSystem : BaseSystem
    {
        // 当前detail角色        
        private Cproto.CharCard _cardDetail;
        List<AttrItem> _attrs = new List<AttrItem>();

        public override string GetName()
        {
            return "charcard";
        }

        public List<AttrItem> GetAttrs()
        {
            return _attrs;
        }

        public CharCardSystem()
        {
            registerCmd<Cproto.CardAttrs>("initattrs", onCmdInitAttrs);
            registerCmd<Cproto.CardAttrs>("attrs", onCmdAttrs);
        }

        public void SetCardDetail(int id)
        {
            _cardDetail = Card.DataCenter.It.GetCard(id);
        }

        public Cproto.CharCard GetCardDetail()
        {
            return _cardDetail;
        }

        protected void onCmdInitAttrs(object data)
        {
            Cproto.CardAttrs args = data as Cproto.CardAttrs;
            Log.LogCenter.Default.Debug($"onCmdInitAttrs: {args.Attrs.Count}");
            var attrs = args.Attrs;

            _attrs.Clear();
            foreach(var item in attrs)
            {
                //Log.LogCenter.Default.Debug($"  attr: {item.Index} value: {item.Value}");
                var one = new AttrItem();
                one.index = item.Index;
                one.value = item.Value;
                _attrs.Add(one);
            }
        }

        public AttrItem GetAttr(int index)
        {
            foreach(var item in _attrs)
            {
                if (item.index == index)
                    return item;
            }
            return null;
        }

        protected void onCmdAttrs(object data)
        {
            Cproto.CardAttrs args = data as Cproto.CardAttrs;
            Log.LogCenter.Default.Debug($"onCmdAttrs: {args.Attrs.Count}");
            var attrs = args.Attrs;
            
            foreach (var one in attrs)
            {
                var item = GetAttr(one.Index);
                if(item == null)
                {
                    // append it
                    item = new AttrItem();
                    item.index = one.Index;
                    item.value = one.Value;
                    _attrs.Add(item);
                } 
                else
                {
                    item.value = one.Value;
                }              
            }
            // events
        }

        public void SetEquip(int index, string id, Action<int> cb)
        {
            if (index < 0 || index >= CardConsts.MaxEquip)
                return;
            Request<Cproto.EmptyArg>("setequip",
                new Cproto.CardSetEquip() {
                    Index= index, Id = id
                },
                (result, code) => {
                    if(code == 0)
                    {
                        var card = _cardDetail;
                        sureEquips(card, index);
                        card.Equips[index].EquipId = id;
                    }
                    cb(code);
            });
        }

        private void sureEquips(Cproto.CharCard card, int index)
        {
            if(index >= card.Equips.Count)
            {
                for(var i = card.Equips.Count; i < CardConsts.MaxEquip; i ++)
                {
                    card.Equips.Add(new Cproto.EquipSlot());
                }
            }            
        }

        public void SetSkill(int index, string id, Action<int> cb)
        {
            if (index < 0 || index >= CardConsts.MaxSkill)
                return;

            Request<Cproto.EmptyArg>("setskill",
                new Cproto.CardSetSkill()
                {
                    Index = index,
                    Id = id
                },
                (result, code) => {
                    if (code == 0)
                    {
                        var card = _cardDetail;

                        sureSkills(card, index);
                        card.Skills[index].SkillId = id;
                        card.Skills[index].Level = 1;
                    }
                    cb(code);
                });
        }

        private void sureSkills(Cproto.CharCard card, int index)
        {
            if (index >= card.Skills.Count)
            {
                for (var i = card.Skills.Count; i < CardConsts.MaxSkill; i++)
                {
                    card.Skills.Add(new Cproto.SkillSlot());
                }
            }
        }

        public void SaveCurCard()
        {
            Request<Cproto.EmptyArg>("savecard", new Cproto.EmptyArg(), (result, code) => { 
            });
        }
    }
}
