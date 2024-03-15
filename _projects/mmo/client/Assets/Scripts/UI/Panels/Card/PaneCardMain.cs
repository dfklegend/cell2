using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Log;
using Network;
using System.Text;

namespace Phoenix.Game
{
    [StringType("PanelCardMain")]
    public class PanelCardMain : BasePanel
    {   
        private Button _btnDo;
        private Button _btnShortcut;
        private Button _btnCards;
        private Text _charInfo;
        private InputField _input;

        ScrollRect _scrollRect;
        private Text _logs;
        public override void OnReady()
        {
            SetDepth(100);
            base.OnReady();

            _btnDo = _root.Find("BG/bottom/btnDo").GetComponent<Button>();
            _btnDo.onClick.AddListener(() => {
                onDo();
            });

            _btnShortcut = _root.Find("BG/bottom/btnShortcut").GetComponent<Button>();
            _btnShortcut.onClick.AddListener(() => {
                onShortcut();
            });

            _btnCards = _root.Find("BG/up/btnCards").GetComponent<Button>();
            _btnCards.onClick.AddListener(() => {
                onCards();
            });

            var btnFight = _root.Find("BG/up/btnFight").GetComponent<Button>();
            btnFight.onClick.AddListener(() => {
                onFight();
            });

            var btnEquipList = _root.Find("BG/up/btnEquipList").GetComponent<Button>();
            btnEquipList.onClick.AddListener(() => {
                onEquipList();
            });

            var btnSkillList = _root.Find("BG/up/btnSkillList").GetComponent<Button>();
            btnSkillList.onClick.AddListener(() => {
                onSkillList();
            });


            _input = _root.Find("BG/bottom/inputCmd").GetComponent<InputField>();


            _charInfo = _root.Find("BG/up/info/bg/info").GetComponent<Text>();

            _scrollRect = TransformUtil.FindComponent<ScrollRect>(_root, "BG/bottom/logs/scroll");
            _logs = _root.Find("BG/bottom/logs/scroll/Text").GetComponent<Text>();

            _logs.text = "";

            refreshCharInfo(Card.DataCenter.It.charInfo);
            BindEvents(true);
        }

        public override void OnDestroy()
        {
            BindEvents(false);
        }

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;            
            events.Bind(Card.EventDefine.CharInfo, onCharInfo, bind);
            events.Bind(Card.EventDefine.BattleLog, onBattleLog, bind);
        }

        private void onCharInfo(params object[] args)
        {
            var e = args[0] as Card.HEventCharInfo;
            refreshCharInfo(e.charInfo);
        }

        private void onBattleLog(params object[] args)
        {
            var e = args[0] as Card.HEventBattleLog;
            addLog(e.log);
        }

        private void onDo()
        {
            string cmd = _input.text;
            if(string.IsNullOrEmpty(cmd))
            {
                addLog("请输入命令");
                return;
            }
            GMCmdMgr.It.Execute(cmd, null);
        }

        private void onShortcut()
        {
            UIMgr.It.GetPanel("PanelShortcut").Show();
        }

        private void onCards()
        {
            UIMgr.It.GetPanel("PanelCards").Show();
        }

        private void onFight()
        {
            UIMgr.It.GetPanel("PanelCardFight").Show();
        }

        private void onEquipList()
        {
            BagCtrl.It.ShowEquipBag();
        }

        private void onSkillList()
        {
            BagCtrl.It.ShowSkillBag();
        }

        private void refreshCharInfo(Cproto.CharInfo info)
        {
            if (_charInfo == null)
                return;
            var sb = new StringBuilder();

            sb.Append($"名字: {info.Name}\n");
            sb.Append($"等级: {info.Level}\n");
            sb.Append($"Exp: {info.Exp}\n");
            sb.Append($"金币: {info.Money}\n");

            _charInfo.text = sb.ToString();
        }

        private void addLog(string log)
        {            
            _logs.text += log + "\n";
            _scrollRect.verticalNormalizedPosition = 0f;
        }
    }
} // namespace Phoenix
