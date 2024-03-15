using System;
using System.Text;

namespace Phoenix.Utils
{
    // 生成数据并校验
    // 随机4位数字一组
    // token就是把数字加起来然后%10000
    public class CorrectnessData
    {
        const int STEP = 4;
        private string _content;
        public string content { get { return _content; } }
        private string _token;
        public string token { get { return _token; } }
        
        Random _rand = new Random(Environment.TickCount/100000);

        public void Generate()
        {
            _rand = new Random(Environment.TickCount / 100000);

            _content = GenerateContent(_rand.Next(10, 200));
            _token = CalcToken(_content);
        }

        public string GenerateContent(int num)
        {
            StringBuilder sb = new StringBuilder();
            for(int i = 0; i < num; i ++)
            {
                for(int j = 0; j < STEP; j ++)
                    sb.Append('1' + _rand.Next(0, 9));
            }
            return sb.ToString();
        }

        public string CalcToken(string content)
        {
            int total = 0;
            var len = content.Length / STEP;
            
            int index = 0;
            for(var i = 0; i < len; i ++)
            {
                string one = "";
                for (var j = 0; j < STEP; j ++)
                {                    
                    one += content[index++];
                }
                total += int.Parse(one);
            }
            return (total%10000).ToString();
        }
    }
} // Phoenix.Utils