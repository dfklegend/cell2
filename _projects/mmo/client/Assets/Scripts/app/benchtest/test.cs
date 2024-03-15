using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Client;
using Phoenix.Log;
using Phoenix.Core;
using Benchtest;
using Phoenix.Scheduler;
using System.Threading;
using System.Threading.Tasks;
using Phoenix.Network;
using Phoenix.Utils;
using Phoenix.Network.Protocol.Pomelo;

public class BenchTest 
{
    const int maxCallTimes = 100000;
    public static ChatClient chatClient = new ChatClient();

    public void Start()
    {   
        chatClient.Start("127.0.0.1", 30021);
        testClient(chatClient).NoWarnning();
    }   


    static async Task testClient(ChatClient client)
    {
        // 每毫秒发送一个消息
        // 一共发送1w个
        while (!client.IsReadyForTest())
        {
            await Task.Delay(1);
        }

        RPCStat stat = new RPCStat();
        int rest = maxCallTimes;
        // speed/ms
        const int speed = 5;
        var lastTime = TimeUtil.NowTick();
        stat.Start();
        while (rest > 0)
        {
            await Task.Delay(1);
            var now = TimeUtil.NowTick();
            int needCall = (now - lastTime) * speed;
            if (needCall > rest)
                needCall = rest;

            for (var i = 0; i < needCall; i++)
            {
                CallHello(client, stat);
            }

            lastTime = now;
            rest -= needCall;

            stat.TryDump();
        }

        // 等待执行完毕
        while (stat.GetTotalTimes() < 0.9 * maxCallTimes)
        {
            await Task.Delay(1);
            stat.TryDump();
        }

        PConsole.Log("Test over!!!!");
        stat.DumpTotal();
    }

    static void CallHello(ChatClient client, RPCStat stat)
    {
        var start = TimeUtil.Now();
        client.Hello((result) => {
            if (result.IsSucc())
            {
                stat.AddOneCost(TimeUtil.Now() - start);
            }
        });
    }
}
