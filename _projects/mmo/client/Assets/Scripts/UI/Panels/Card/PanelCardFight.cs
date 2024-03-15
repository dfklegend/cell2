using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Log;
using System.Collections.Generic;

namespace Phoenix.Game
{
    [StringType("PanelCardFight")]
    public class PanelCardFight : BasePanel
    {
        private int _idUp = -1;
        private int _idDown = -1;
        private Text _infoUp;
        private Text _infoDown;

        

        public override void OnReady()
        {
            SetDepth(PanelDepth.Normal-1);
            base.OnReady();

            var btnCancel = _root.Find("BG/btnCancel").GetComponent<Button>();
            btnCancel.onClick.AddListener(() => {
                onCancel();
            });

            var btnStart = _root.Find("BG/btnStart").GetComponent<Button>();
            btnStart.onClick.AddListener(() => {
                onStart();
            });

            var btnUpSelect = _root.Find("BG/up/btnSelect").GetComponent<Button>();
            btnUpSelect.onClick.AddListener(() => {
                onUpSelect();
            });
            _infoUp = _root.Find("BG/up/bg/Text").GetComponent<Text>();

            var btnDownSelect = _root.Find("BG/down/btnSelect").GetComponent<Button>();
            btnDownSelect.onClick.AddListener(() => {
                onDownSelect();
            });
            _infoDown = _root.Find("BG/down/bg/Text").GetComponent<Text>();
        }

        private void clearSelect()
        {
            _idUp = -1;
            _idDown = -1;
            setInfo(_infoUp, null);
            setInfo(_infoDown, null);
        }

        private void onCancel()
        {
            clearSelect();
            Hide();
        }

        private void onStart()
        {            
            // 检查是否有选择目标
            if(_idUp == -1 || _idDown == -1)
            {
                UIMgr.It.GetPanel<PanelDialog>().ShowInfo("请选择角色");
                return;
            }

            ClientApp.It.client.StartFight(_idDown, _idUp);
            Hide();
        }

        private List<Cproto.CharCard> makeDataSources()
        {
            List<Cproto.CharCard> cards = new List<Cproto.CharCard>();
            var from = Card.DataCenter.It.cards;
            for(var i = 0; i < from.Count; i ++)
            {
                var c = from[i];
                if (c.Id == _idUp || c.Id == _idDown)
                    continue;
                cards.Add(c);
            }
            return cards;
        }

        private void onUpSelect()
        {
            var source = makeDataSources();
            UIMgr.It.GetPanel<PanelCards>().SelectFrom(source, (card) => {
                onSelectResult(ref _idUp, _infoUp, card);
            });
        }

        private void onDownSelect()
        {
            var source = makeDataSources();
            UIMgr.It.GetPanel<PanelCards>().SelectFrom(source, (card) => {
                onSelectResult(ref _idDown, _infoDown, card);
            });
        }

        private void onSelectResult(ref int id, Text info, Cproto.CharCard card)
        {
            if (card != null)
            {
                id = card.Id;
            }
            else
            {
                id = -1;
            }
            setInfo(info, card);
        }

        private void setInfo(Text info, Cproto.CharCard card)
        {
            if(card == null)
            {
                info.text = "";
                return;
            }

            info.text = CardItem.MakeInfo(card);
        }
    }
} // namespace Phoenix
