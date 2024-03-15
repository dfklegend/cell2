using System;

namespace Phoenix.Utils
{
    // 非阻塞读取
    public class ConsoleCmdReader
    {
        private string _cmd = "";
        private string _buf = "";
        public void Update()
        {
            if (!Console.KeyAvailable)
                return;
            var key = Console.ReadKey();
            if(key.Key == ConsoleKey.Enter)
            {
                _cmd = _buf;
                _buf = "";
                Console.WriteLine("");
                return;
            }

            _buf += key.KeyChar;
        }

        public bool HasCmd()
        {
            return !string.IsNullOrEmpty(_cmd);
        }

        public string PopCmd()
        {
            var str = _cmd;
            _cmd = "";
            return str;
        }
    }
}
