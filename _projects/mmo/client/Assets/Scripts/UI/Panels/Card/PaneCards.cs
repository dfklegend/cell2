using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using UnityEngine.Events;
using Phoenix.Log;
using Network;
using System.Text;
using System.Collections.Generic;
using System;

namespace Phoenix.Game
{

    public enum ePanelCardsState
    {
        Normal,     // 普通
        Delete,     // 删除
        Select,     // 选择
    }

    public class CardItem
    {
        Transform _root;
        RectTransform _rect;
        Button _btn;
        Text _text;
        Cproto.CharCard _data;

        public void Init(Transform tran)
        {
            _root = tran;
            _btn = tran.GetComponent<Button>();
            _text = tran.Find("Text").GetComponent<Text>();
            _rect = tran.GetComponent<RectTransform>();
        }

        public Button GetBtn()
        {
            return _btn;
        }

        public Cproto.CharCard GetData()
        {
            return _data;
        }

        public void SetInfo(Cproto.CharCard card)
        {
            _data = card;           
            setText(MakeInfo(card));
        }

        public static string MakeInfo(Cproto.CharCard card)
        {
            var sb = new StringBuilder();
            sb.AppendFormat("id: {0}\n", card.Id);
            sb.AppendFormat("name: {0}\n", card.Name);
            sb.AppendFormat("战绩: {0}/{1}\n", card.Stat.Win, card.Stat.Total);
            return sb.ToString();
        }


        private void setText(string str)
        {
            _text.text = str;
        }

        public void Show(bool show)
        {
            _root.gameObject.SetActive(show);
        }

        public void Destroy()
        {
            GameObject.Destroy(_root.gameObject);
        }

        public void SetPos(float x, float y)
        {
            _rect.anchoredPosition = new Vector2(x, y);
        }
    }

    [StringType("PanelCards")]
    public class PanelCards : BasePanel
    {
        private const int BtnBeginX = -180;
        private const int BtnBeginY = 320;
        private const int BtnWidth = 180;
        private const int BtnHeight = 100;
        private const int ColNum = 3;

        private Button _btnClose;
        private GameObject _cardPrefab;
        private Text _info;
        
        private List<CardItem> _cards = new List<CardItem>();

        private List<Cproto.CharCard> _source = null;
        private Action<Cproto.CharCard> _cbSelect = null;


        private Transform _normal;
        private Transform _select;        
        private ePanelCardsState _state = ePanelCardsState.Normal;

      
        public override void OnReady()
        {
            SetDepth(PanelDepth.Normal);
            base.OnReady();

            _btnClose = _root.Find("BG/normal/btnClose").GetComponent<Button>();
            _btnClose.onClick.AddListener(() => {
                Hide();
            });

            var btnCreate = _root.Find("BG/normal/btnCreate").GetComponent<Button>();
            btnCreate.onClick.AddListener(() => {
                onCreateCard();
            });

            var btnDelete = _root.Find("BG/normal/btnDelete").GetComponent<Button>();
            btnDelete.onClick.AddListener(() => {
                toDeleting();
            });

            var btnDeleteCancel = _root.Find("BG/select/btnCancel").GetComponent<Button>();
            btnDeleteCancel.onClick.AddListener(() => {
                onCancelSelect();
            });

            _cardPrefab = _root.Find("BG/cards/prefab").gameObject;
            _cardPrefab.SetActive(false);
            _normal = _root.Find("BG/normal");
            _select = _root.Find("BG/select");
            _info = _root.Find("BG/select/info").GetComponent<Text>();

            toNormal();
            BindEvents(true);
        }

        public override void OnDestroy()
        {
            BindEvents(false);
        }

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;
            events.Bind(Card.EventDefine.RefreshCards, onShowCards, bind);
        }

        protected override void onShow()
        {
            base.onShow();
            showCards();
        }

        protected override void onHide()
        {
            base.onHide();
            destroyCards();
            _source = null;
        }

        private void onShowCards(params object[] args) 
        {
            if (!IsVisible())
                return;
            showCards();
        }

        private List<Cproto.CharCard> GetDataSource()
        {
            if (_source != null)
                return _source;
            return Card.DataCenter.It.cards;
        }

        private void showCards()
        {
            destroyCards();
            var data = GetDataSource(); ;
            if (data.Count == 0)
                return;
            var index = 0;
            for(var row = 0; row < 1+data.Count/ColNum; row ++)
            {
                for (var col = 0; col < ColNum; col++, index ++)
                {
                    if (index >= data.Count)
                        break;

                    var go = GameObject.Instantiate(_cardPrefab);
                    go.transform.SetParent(_cardPrefab.transform.parent, false);
                    go.name = "card" + index;
                    var item = new CardItem();
                    item.Init(go.transform);
                    item.SetInfo(data[index]);
                    item.SetPos(BtnBeginX + col * BtnWidth, BtnBeginY - row * BtnHeight);

                    item.GetBtn().onClick.AddListener(createOnClickAction(item));

                    item.Show(true);
                    _cards.Add(item);           
                }
            }            
        }        

        private void destroyCards()
        {
            for(var i = 0; i < _cards.Count; i ++)
            {
                _cards[i].Destroy();
            }
            _cards.Clear();
        }

        private UnityAction createOnClickAction(CardItem item)
        {
            return () => 
            {
                onClickItem(item);
            };
        }

        private void onCreateCard()
        {
            UIMgr.It.GetPanel<PanelCreateCard>().Show();
        }

        private void showSelect()
        {
            _normal.gameObject.SetActive(false);
            _select.gameObject.SetActive(true);
        }

        private void toDeleting()
        {
            setInfo("点击删除卡牌");
            _state = ePanelCardsState.Delete;
            showSelect();
        }

        private void toSelect()
        {
            setInfo("点击选择");
            _state = ePanelCardsState.Select;
            showSelect();
        }

        private void toNormal()
        {    
            _state = ePanelCardsState.Normal;
            _normal.gameObject.SetActive(true);
            _select.gameObject.SetActive(false);
        }

        private void setInfo(string str)
        {
            _info.text = str;
        }

        private void onCancelSelect()
        {
            if (_state == ePanelCardsState.Delete)
            {
                toNormal();
                return;
            }
            if (_state == ePanelCardsState.Select)
            {
                if (_cbSelect != null)
                    _cbSelect(null);
                toNormal();
                Hide();
            }
        }

        private void onClickItem(CardItem item)
        {            
            if (_state == ePanelCardsState.Delete)
            {
                GMCmdMgr.It.Execute($"lcmd deletecard {item.GetData().Id}", null);
                return;
            }
                
            if(_state == ePanelCardsState.Select)
            {
                if (_cbSelect != null)
                    _cbSelect(item.GetData());
                toNormal();
                Hide();
                return;
            }

            // show detail
            GMCmdMgr.It.Execute($"lcmd opencard {item.GetData().Id}", (succ)=> {
                if(succ)
                {
                    Systems.It.GetSystem<CharCardSystem>("charcard").SetCardDetail(item.GetData().Id);
                    UIMgr.It.GetPanel<PanelCardDetail>().Show();
                }                    
            });
            
        }

        public void SelectFrom(List<Cproto.CharCard> source, Action<Cproto.CharCard> cb)
        {
            if (!IsReady())
                return;
            _source = source;
            _cbSelect = cb;
            toSelect();
            Show();
        }
        
    }
} // namespace Phoenix
