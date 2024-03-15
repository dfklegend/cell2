using Phoenix.Network;
using Phoenix.Scheduler;
using System;
using System.Net;
using System.Runtime.CompilerServices;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Phoenix.Network.Protocol.Pomelo;
using SimpleJson;
using Phoenix.Utils;
using PomeloCommon;
using Benchtest;

public class MainArgs
{
    public int clientNum = 1;
    public int msgSize = 1024;
}

class Program
{
    static NetConfig env;
    static MainArgs mainArgs = new MainArgs();
    static string content = "";
    public static ChatClient chatClient = new ChatClient();

    static void Main(string[] args)
    {
        Console.WriteLine("\nApp.start\n");
        
        ThreadMgr.Start(2);               
        
        var client = chatClient;
        chatClient.Start("127.0.0.1", 30021);
        
        Console.WriteLine("Begin loop");
        Console.WriteLine("Input exit to exit app");
        ConsoleCmdReader cmdBuf = new ConsoleCmdReader();
        while (true)
        {
            cmdBuf.Update();
            if(cmdBuf.HasCmd())
            {
                var cmd = cmdBuf.PopCmd();
                if(cmd == "exit")
                {
                    break;
                }
                processCmd(cmd);
            }          
            ThreadMgr.RunStep();            
        

            Thread.Sleep(1);
        }
   
        ThreadMgr.Stop();

        Console.WriteLine("\napp exit\n");
    }  

    static void processCmd(string cmd)
    {
        
    }
}
