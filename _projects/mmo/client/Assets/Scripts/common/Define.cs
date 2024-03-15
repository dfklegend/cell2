using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace PomeloCommon
{
    public class Interest
    {
        public string name = "兴趣";
        public int exp = 100;
    }

    public class Person
    {
        public int serialId = 0;
        public string name = "haha";
        public int age = 30;        
        public List<string> tags = new List<string>();
        public List<Interest> interests = new List<Interest>();
        public string content = "";
    }

    public class PersonTitle
    {
        public int serialId = 0;
        public string name = "h";
        public string title = "title";
    }

    public class Define
    {
        public static Person MakePerson()
        {
            var p1 = new Person();
            p1.tags.Add("帅哥");
            p1.tags.Add("有钱");

            p1.interests.Add(new Interest { name = "游戏" });
            p1.interests.Add(new Interest { name = "篮球" });
            return p1;
        }
    }
}
