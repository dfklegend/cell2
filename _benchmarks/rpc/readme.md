
Pomelo 对比测试
    https://github.com/NetEase/pomelo/wiki/pomelo-rpc%E6%80%A7%E8%83%BD%E6%B5%8B%E8%AF%95%E6%8A%A5%E5%91%8A
    
    数据摘要
    (connector和echo存在request转发)
    场景A
    connector和echo业务进程各1个.
    2个客户端并发, 每隔1ms发起一次request请求(msg='Hello World'), 每个客户端总计发送1w次, 服务器对每个request回复一个200.
    服务器完成2w次请求的时间为14.835s, 平均1348次/s.
    服务器完成一次RPC调用的时间约为: 2~8ms
    在服务器运行过程中: connector进程对CPU的占用平均值为: 91.6% [CPU占用的采样点为: 92%, 94%, 95%, 87%, 84%, 96%, 93%]; echo进程对CPU的占用平均值为: 28.1% [CPU占用的采样点为: 30%, 20%, 33%, 22%, 25%, 46%, 21%]
    在客户端运行过程中: client进程对CPU的占用平均值为: 30.1% [CPU占用的采样点为: 18%, 24%, 25%, 40%, 16%, 49%, 39%]
    场景B
    4个connector和1个echo业务进程.
    4个客户端并发且分别连接1个connector, 每隔1ms发起一次request请求(msg='Hello World'), 每个客户端总计发送1w次, 服务器对每个request回复一个200.
    服务器完成4w次请求的时间为14.866s, 平均2690次/s.
    服务器完成一次RPC调用的时间约为: 1~25ms
    在服务器运行过程中: connector进程对CPU的占用平均值为: 71.8% [CPU占用的采样点为: 75%, 71%, 71%, 74%, 68%]; echo进程对CPU的占用平均值为: 81.3% [CPU占用的采样点为: 81%, 82%, 83%, 79%]
    在客户端运行过程中: client进程对CPU的占用平均值为: 28.0% [CPU占用的采样点为: 28%, 29%, 29%, 26%]


    测试1
    单进程 all-1
    发送速度 5req /s
    avg 0.2ms /req  
    net send 22MB/2

    发送速度 10req /s
    avg speed 0.13 ms /req  
    net send 40MB/2

    发送速度 20req /s
    avg speed 0.13 ms /req  
    net send 80MB/2
    明显看出来response延迟处理

    测试2
    本机多进程
    gate-1
    chat-1

    发送速度 10req /s
    avg speed 0.14 ms /req   

    测试3
    注: 去掉服务器无用的log输出后(每个消息都有)，提升很大
    单进程
    单客户端
    发送速度 100req /s
    avg speed 0.01 ms /req
    基本正常, 100000w req/s

    发送速度 200req /s
    avg speed 0.008 ms /req
    但是看起来是先集中发了很多命令，再收
    应该接近极限
    
    上面主要是统计总体收到多少，下面统计具体req的响应时间
    统计每一个req的具体返回时间
    localhost测试

    发送速度: 1 req/ms
    avg response time: 0.71 ms
    RPS: 998 req/s

    发送速度: 100 req/ms
    avg response time: 10.81 ms
    RPS: 80000 req/s

    发送速度: 150 req/ms
    avg response time: 104.52 ms 
    RPS: 100000 req/s
    (明显有卡顿情况)


    多连接测试
    本机    
    发送速度: 1req /ms

    gate x1 
    chat x1
    并发: 1000 
    消息超时

    并发: 100
    avg response time: 54-600 ms
    每连接 RPS: 800 req/s

    并发: 50
    avg response time: 6 ms
    每连接 RPS: 996 req/s

    50属于性能合适，没有恶化 后面用这个指标，多进程测试
    
    gate x2 
    chat x1

    并发: 50
    avg response time: 9-10 ms
    每连接 RPS: 999 req/s

    gate x2 
    chat x2

    并发: 50
    avg response time: 6-7 ms
    每连接 RPS: 998 req/s

    基本结论
    总体吞吐量在80000req/s左右    
    性能恶化之后，如何保证服务器状况需要考虑
    (C#实现大概 10000req/s)
    (Cell大概 5000-10000 req/s)

    ----
    sessionData调整为 被动请求后
    对比测试

    本机    
    发送速度: 1req /ms

    gate x1 
    chat x1    

    并发: 100
    avg response time: 10 ms
    每连接 RPS: 999 req/s

    并发: 150
    avg response time: 20 ms
    每连接 RPS: 999 req/s

    对比结果
    吞吐量又得到提升，主要是减少了sessionData的 json压缩解压，释放了cpu

    ----
    限帧 超过处理maxcost后 sleep时间固定 1ms

    并发: 50
    avg response time:  8 ms
    每连接 RPS:   997 req/s

    并发: 100
    avg response time:  10 ms
    每连接 RPS:   999 req/s

    并发: 150
    avg response time:  30 ms
    每连接 RPS:  999 req/s

    并发: 200
    avg response time:  6000 ms
    每连接 RPS:  635 req/s
    (chat消息队列 > 100000)

    多开一个chat-2
    并发: 200
    avg response time:  40-60 ms
    每连接 RPS:   930 req/s
    (chat消息队列正常，之前测试服务器瓶颈在chat)

    ----
    2023/01/11 复测
    并发: 50
    avg response time :4-9 ms
    999 req/s

    并发: 100
    avg response time:  
    每连接 RPS:   900 req/s
    产生了很多timeout，无法正常工作

    (感觉差了很多，也可能是之前测试并不是那么健壮，并没有发现timeout)
    主要是chat处理不过来，消息堆积太多

    3个chat
    并发: 100
    avg response time: 60-90  
    每连接 RPS:   895 req/s
    总吞吐900x100=90000 req/s

    优化
    ----
    3个chat
    并发: 100
    handlerContext使用pool后
    (之前不使用pool，有bug)
    avg response time: 40-180  
    每连接 RPS:   870 req/s

    ----
    3个chat
    并发: 100
    去掉pool，使用共享变量，同时backsession也改成了共享变量
    avg response time: 20-30 
    每连接 RPS:   960 req/s
    
    ----
    handlercontext和backsession全部创建
    3个chat
    并发: 100
    avg response time: 20-60
    每连接 RPS: 971  
    总吞吐950x100= 95000 req/s