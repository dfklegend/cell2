using System.Text;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    [StringType("PanelFightFE")]
    public class PanelFightFE : BasePanel
    {
        Transform _charRoot;
        Text _logs;
        int _logNum = 0;
        ScrollRect _scrollRect;

        Text _leftContent;
        Text _rightContent;

        Button _btnSelectEnemy;
        Button _btnSelectRole;
        Button _btnRestart;

        Text _sceneInfo;
        Text _posInfo;
        private Text _charInfo;

        public override void OnReady()
        {
            SetDepth(101);
            base.OnReady();

            _charRoot = _root.Find("BG/chars");
            _scrollRect = TransformUtil.FindComponent<ScrollRect>(_root, "BG/logs/scroll");
            _logs = TransformUtil.FindComponent<Text>(_root, "BG/logs/scroll/Text");
            _leftContent = TransformUtil.FindComponent<Text>(_root, "BG/left/scroll/Text");
            _rightContent = TransformUtil.FindComponent<Text>(_root, "BG/right/scroll/Text");
            _btnSelectEnemy = TransformUtil.FindComponent<Button>(_root, "BG/btnSelectEnemy");
            _btnSelectEnemy.onClick.AddListener(OnSelectEnemy);
            _btnSelectEnemy = TransformUtil.FindComponent<Button>(_root, "BG/btnSelectRole");
            _btnSelectEnemy.onClick.AddListener(OnSelectRole);
            _btnRestart = TransformUtil.FindComponent<Button>(_root, "BG/btnRestart");
            _btnRestart.onClick.AddListener(OnRestart);

            var btnBag = TransformUtil.FindComponent<Button>(_root, "BG/btnBag");
            btnBag.onClick.AddListener(OnBag);

            var btnExit = _root.Find("BG/righttop/btnExit").GetComponent<Button>();
            btnExit.onClick.AddListener(() => {
                GameLogicCtrl.It.Clear();

                ClientApp.It.stateCtrl.ChangeState((int)eAppState.CardMain);

                Destroy();
            });
            _sceneInfo = TransformUtil.FindComponent<Text>(_root, "BG/righttop/sceneInfo");            
            _posInfo = TransformUtil.FindComponent<Text>(_root, "BG/righttop/posInfo");
            _charInfo = TransformUtil.FindComponent<Text>(_root,"BG/up/info/bg/info");
            

            var btnGraph = _root.Find("BG/btnGraph").GetComponent<Button>();
            btnGraph.onClick.AddListener(() => {
                UIMgr.It.GetPanel<PanelFightGraph>().Show();
            });

            var btnCmds = _root.Find("BG/btnCmds").GetComponent<Button>();
            btnCmds.onClick.AddListener(this.OnCmds);

            var btnDaySign = _root.Find("BG/btnDaySign").GetComponent<Button>();
            btnDaySign.onClick.AddListener(this.onDaySign);

            _logs.text = "";

            BindEvents(true);
            //Timer
        }

        public override void OnDestroy()
        {
            BindEvents(false);
        }

        public Transform GetCharRoot()
        {
            return _charRoot;
        }

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;
            events.Bind(Card.EventDefine.Attack, OnAttack, bind);
            events.Bind(Card.EventDefine.StartSkill, OnStartSkill, bind);
            events.Bind(Card.EventDefine.InitFight, OnInitFight, bind);
            events.Bind(Card.EventDefine.BattleLog, onBattleLog, bind);
            events.Bind(Card.EventDefine.RefreshSceneInfo, onRefreshSceneInfo, bind);
            events.Bind(Card.EventDefine.CharInfo, onCharInfo, bind);
        }

        protected override void onShow() 
        {
            refreshSceneInfo();
            refreshCharInfo(Card.DataCenter.It.charInfo);
        }

        private string buildAttackLog(Card.HEventAttack e)
        {
            var sd = Card.SkillDataMgr.It.GetItem(e.skillId);
            if (sd == null)
                return "";

            string skillName = sd.name;            

            if (!e.result.data.hit)
            {
                return $"{e.src.name}使用{skillName}攻击{e.tar.name}，未命中\n";
            }
            if (e.result.data.critical)
            {
                return $"{e.src.name}使用{skillName}攻击{e.tar.name}，暴击造成{(int)e.result.data.Dmg}点伤害\n";
            }
            return $"{e.src.name}使用{skillName}攻击{e.tar.name}，造成{(int)e.result.data.Dmg}点伤害\n";
        }
        private string buildWeapon(Card.HEventAttack e)
        {
            string weapon = "主手武器";
            if (e.result.data.hand == Skill.eHandType.OffHand)
                weapon = "副手武器";

            if (!e.result.data.hit)
            {
                return $"{e.src.name}使用{weapon}攻击{e.tar.name}，未命中\n";
            }
            if (e.result.data.critical)
            {
                return $"{e.src.name}使用{weapon}攻击{e.tar.name}，暴击造成{(int)e.result.data.Dmg}点伤害\n";
            }
            return $"{e.src.name}使用{weapon}攻击{e.tar.name}，造成{(int)e.result.data.Dmg}点伤害\n";
        }


        // 只显示目标的伤害，和本方的有益
        private void OnAttack(params object[] args)
        {
            var e = args[0] as Card.HEventAttack;

            // TODO: 控制长度
            var log = buildAttackLog(e);
            addLog(log);

            Log.LogCenter.Default.Debug(log);
            //int selfId = 1;

            // 找到目标的坐标
            var tarModel = Card.UnitModelMgr.It.GetModel(e.tar.id);
            //var srcModel = UnitModelMgr.It.GetModel(e.src.id);
            if (tarModel == null)
                return;
            if (e.result.data.hit)
            {
                if (true/*e.src.id == selfId*/)
                {
                    string str;
                    Color c = Color.white;
                    if (e.result.data.critical)
                    {
                        str = $"暴击{(int)e.result.data.Dmg}";
                        c = Color.yellow;
                    }
                    else
                        str = $"{(int)e.result.data.Dmg}";
                    Vector2 pos = UIUtil.WorldToScreen(null, tarModel.GetWorldPos());
                    var hitText = FText.FloatTextCreator.CreateLineUp(pos,
                        str, 3f, 100f, c);
                    if (e.result.data.critical)
                    {
                        hitText.PlayAnim("critical");
                        tarModel.PlayAnim("behit", 0.5f);
                    }
                }


                if (e.result.data.block/* && e.tar.id == selfId*/)
                {
                    var one = FText.FloatTextCreator.CreateLineUp(
                        UIUtil.WorldToScreen(null, tarModel.GetWorldPos()),
                        "格挡", 3f, 100f, Color.blue);
                    one.OffsetPos(new Vector2(70, 0));
                }
                //tarModel.PlayAnim("behit", 0.5f);
            }
            else
            {
                if (true/*e.tar.id == selfId*/)
                {
                    Vector2 pos = UIUtil.WorldToScreen(null, tarModel.GetWorldPos());
                    var one = FText.FloatTextCreator.CreateLineUp(pos,
                        $"闪避", 3f, 100f, Color.yellow);
                    one.OffsetPos(new Vector2(70, 0));
                }

            }
        }

        private void OnStartSkill(params object[] args)
        {
            var e = args[0] as Card.HEventStartSkill;

            var srcModel = Card.UnitModelMgr.It.GetModel(e.src.id);
            var sd = Card.SkillDataMgr.It.GetItem(e.skillId);
            if (sd != null && sd.normalAttack == 0)
            {
                Vector2 pos = UIUtil.WorldToScreen(null, srcModel.GetWorldPos());
                var one = FText.FloatTextCreator.CreateLineUp(pos,
                    $"{sd.name}", 3f, 100f, Color.yellow);
                one.OffsetPos(new Vector2(70, 0));
            }
        }

        private void OnInitFight(params object[] args)
        {
            //var left = FightCtrl.It.GetChar(1);

            //_leftContent.text = left.MakeBrief();

            //var right = FightCtrl.It.GetChar(2);
            //_rightContent.text = right.MakeBrief();
        }

        private void OnSelectEnemy()
        {
            //UIMgr.It.OpenPanel("PanelSelectEnemy").Show();
            //UIMgr.It.GetPanel<PanelSelectEnemy>("PanelSelectEnemy").ReqSelect(
            //    FightSimulateCtrl.It.enemyCfgId);
        }

        private void OnSelectRole()
        {
            UIMgr.It.OpenPanel("PanelSelectRole").Show();
        }

        private void OnRestart()
        {
            //FightSimulateCtrl.It.Restart();
        }

        private void OnBag()
        {
            UIMgr.It.OpenPanel("PanelItemBag").Show();
        }

        private void OnCmds()
        {
            UIMgr.It.OpenPanel("PanelShortcut").Show();
        }

        private void onBattleLog(params object[] args)
        {
            var e = args[0] as Card.HEventBattleLog;
            addLog(e.log);
        }

        private void onRefreshSceneInfo(params object[] args)
        {
            var e = args[0] as Card.HEventBattleLog;
            refreshSceneInfo();
        }

        private void refreshSceneInfo()
        {
            string str = $"场景: {Card.DataCenter.It.SceneServer}-{Card.DataCenter.It.SceneId}-{Card.DataCenter.It.SceneCfgId}";
            _sceneInfo.text = str;
        }

        private void addLog(string log)
        {
            _logNum++;
            if (_logNum > 100)
            {
                // just clear
                // TODO: 做保留
                _logs.text = "";
                _logNum = 0;
            }
            _logs.text += log;
            _scrollRect.verticalNormalizedPosition = 0f;
        }

        public override void Update()
        {
            refreshPos();
        }

        private void refreshPos()
        {
            var player = Game.Card.FightCtrl.It.player;
            if (player == null)
                return;
            var pos = player.pos;
            _posInfo.text = $"({pos.x.ToString("F2")},{pos.z.ToString("F2")})";
        }

        private void onDaySign()
        {            
            UIMgr.It.OpenPanel("PanelDaySign").Show();
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

        private void onCharInfo(params object[] args)
        {
            var e = args[0] as Card.HEventCharInfo;
            refreshCharInfo(e.charInfo);
        }
    }
} // namespace Phoenix
