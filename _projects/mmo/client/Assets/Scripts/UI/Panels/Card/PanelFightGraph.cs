using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Log;
using System.Collections.Generic;
using XCharts.Runtime;

namespace Phoenix.Game
{
    [StringType("PanelFightGraph")]
    public class PanelFightGraph : BasePanel
    {
        const int SampleNum = 30;

        GameObject _graph;
        public override void OnReady()
        {
            SetDepth(PanelDepth.AboveNormal + 1);
            base.OnReady();

            var btnClose = _root.Find("BG/btnClose").GetComponent<Button>();
            btnClose.onClick.AddListener(() => {
                onClose();
            });

            var btnRefresh = _root.Find("BG/btnRefresh").GetComponent<Button>();
            btnRefresh.onClick.AddListener(() => {
                onRefresh();
            });
            _graph = _root.Find("BG/graph").gameObject;
            tryInitData(_graph);
        }
        

        private void onClose()
        {            
            Hide();
        }

        protected override void onShow()
        {
            base.onShow();
            tryInitData(_graph);
        }

        private void onRefresh()
        {
            tryInitData(_graph);
        }

        string GetName(bool downner)
        {
            var data = Card.DataCenter.It.fightStat;
            if (data.GetUnitNum() == 0)
            {
                return "";
            }
            var units = data.GetUnits();
            var entityDown = Card.FightCtrl.It.GetEntityDown();

            foreach (var id in units.Keys)
            {
                if (id == entityDown && downner)
                    return units[id].name;
                if(id != entityDown && !downner)
                    return units[id].name;
            }
            return "noset";
        }
        
        void tryInitData(GameObject root)
        {
            var data = Card.DataCenter.It.fightStat;
            if(data.GetUnitNum() == 0 )
            {
                return;
            }

            var chart = root.GetComponent<BaseChart>();
            if (chart == null)
            {
                chart = root.AddComponent<LineChart>();
                chart.Init();
                chart.SetSize(600, 600);
            }
            var title = chart.GetOrAddChartComponent<Title>();
            title.text = "ÉËº¦ÇúÏß";
            title.subText = string.Format("({0})À¶/({1})ÂÌ", GetName(false), GetName(true));

            var tooltip = chart.GetOrAddChartComponent<Tooltip>();
            tooltip.show = true;

            var legend = chart.GetOrAddChartComponent<Legend>();
            legend.show = false;

            var xAxis = chart.GetOrAddChartComponent<XAxis>();
            xAxis.splitNumber = 10;
            xAxis.boundaryGap = true;
            xAxis.type = Axis.AxisType.Category;

            var yAxis = chart.GetOrAddChartComponent<YAxis>();
            yAxis.type = Axis.AxisType.Value;

            chart.RemoveData();            
            chart.AddSerie<Line>();
            chart.AddSerie<Line>();

            var samples = data.GetSamples();
            for (int i = 0; i < SampleNum; i++)
            {
                chart.AddXAxisData("x" + i);                                
            }

            var sampleNum = data.GetSampleNum();
            float factor = (float)sampleNum / SampleNum;
            int index = 0;
            foreach(var id in samples.Keys)
            {
                // 0:À¶ 1:ÂÌ
                const int indexBlue = 0;
                const int indexGreen = 1;
                
                index = id == Card.FightCtrl.It.GetEntityDown() ? indexGreen : indexBlue;
                var itemData = samples[id].data;
                if(sampleNum <= SampleNum)
                {
                    for (int i = 0; i < sampleNum; i++)
                    {
                        chart.AddData(index, itemData[i].dmg);
                    }
                }
                else
                {
                    for (int i = 0; i < SampleNum-1; i++)
                    {
                        int at = (int)((float)i * factor);
                        chart.AddData(index, itemData[at].dmg);
                    }
                    int finalAt = itemData.Count - 1;
                    chart.AddData(SampleNum - 1, itemData[finalAt].dmg);
                }                
            }
        }
    }
} // namespace Phoenix
