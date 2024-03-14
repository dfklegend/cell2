# Cell2服务器架构设计
    基于proto.actor
    RPC base from protobuf(gRPC)
    Service
        特殊actor,将会自动注册到服务发现节点
    Service-discovery
        方便的组建服务器组
    Service-Route
        方便的service路由机制
    Request/Response
        提供callback方式的调用，便于书写现场代码
    Client message auto forward
        客户端消息，基于路由规则，自动投递到目标Service

    优势
    . 很容易基于Service组织服务器组
    . Service可以集成在一个进程内，方便调试
    . Service可以灵活的配置在不同节点
    . 节点伸缩很方便
    . 基于protobuf的协议也很容易扩展外部服务
    . 基于protobuf的协议提供了服务之间的接口兼容性，
      使热更和灰度更新成为可能    

# TODO
    . 错误处理修改
        如何将错误正确的传递给客户端
        能方便的屏蔽协议(客户端显示: 功能暂时关闭)
    . 目录结构调整 pomelo应该只是一个客户端协议实现
    . 客户端协议更方便定制
	. 更好的性能分析
    . 单routine的高效timer模块
	. 进一步搭建模拟MMO工程，验证合适的模型    
    

# 模块说明
## 启动流程
    首先通过节点id来启动某个节点配置
    节点配置中定义了 StartMode: allinone
    
    代码通过定义lanuchmode来约定启动那些module
    baseapp.LaunchFunc("allinone", func(app interfaces.IApp) {
        log.Printf("allinone mode\n")
        utils.NodeAddCommonModules(app)
        utils.NodeAddClusterModules(app)
    })

    同时节点中定义了启动那些服务
    比如:
    all-1:
        StartMode: allinone
        Address: 127.0.0.1:3001
        Services:
            - gate-1
            - gate-2
            - chat-1
            - chat-2

    代表了使用allinone作为launchmode
    使用127.0.0.1:3001作为邮箱地址
    启动gate-1,gate-2,chat-1,chat-2服务    

    具体的启动步骤，可以依次添加module, module将依次启动
    参看 baseapp

    可以使用nodebuilder内NewBuilder()的方式来启动服务器

## 服务模块
    nodeservice模块是服务器功能组合的基本元素
    nodeservice在固定routine执行, 提供基本的timer和event支持
    在nodes.yaml中，配置某个节点启动那些服务
        
    具体服务，配置自己服务类型和是否前端服务
    比如:
    services:
        gate-1:
            Type: gate
            Frontend: true
            WSClientAddress: 127.0.0.1:30011
            ClientAddress: 127.0.0.1:30021
        gate-2:
            Type: gate
            Frontend: true
            WSClientAddress: 127.0.0.1:30012
            ClientAddress: 127.0.0.1:30022
        chat-1:
            Type: chat
        chat-2:
            Type: chat
        chat-3:
            Type: chat
    任何服务目前可以都可以配置成前端服务，目前只支持pomelo协议(tcp和ws)
    如上面配置, gate-1,gate-2都监听了ws和tcp, chat则是后端服务

    服务构建器，前后端服务，需要不同创建方式

### 注册前端服务
    定义服务构建器
    func NewCreator() service.IServiceCreator {
	    return service.NewFuncCreator(func(name string) {
		builder.StartFrontService(name, func() actor.Actor { return NewService() })
	    })
    }
    在AppBuilder中注册服务
        .RegisterServiceCreators(func(factory service.IServiceFactory) {
			factory.Register("gate", gate.NewCreator())
		})

### 注册后端服务
    定义服务构建器
    func NewCreator() service.IServiceCreator {
	    return service.NewFuncCreator(func(name string) {
		builder.StartBackService(name, func() actor.Actor { return NewService() })
	    })
    }
    在AppBuilder中注册服务
        .RegisterServiceCreators(func(factory service.IServiceFactory) {
			factory.Register("logic", logic.NewCreator())
		})

## 前端服务/后端服务
    参看pomelo
    https://github.com/NetEase/pomelo/wiki/%E4%B8%8E%E5%AE%A2%E6%88%B7%E7%AB%AF%E9%80%9A%E4%BF%A1
    所有可以接受客户端连接的，我们称之为前端服务，客户端的消息是直接发送给前端服务
    后端服务就是客户端非直接连接的服务
    client可以向前端和后端投递消息，接口一致
    比如client.request("serviceType.group.func", ...)
    意思调用服务函数，serviceType上的group.func函数
    前端消息，收到时候，如果发现当前前端服务类型就是serviceType的服务，那么就直接处理了
    如果不是，那么要做转发，后端服务的消息，会根据routeFunc(frontSession)来计算转发的目标服务，并转发

## 向客户端推送消息
    可以调用
    app.PushMessageById向客户端主动推送消息

## 客户端消息顺序
    客户端收到的消息顺序严格按照代码顺序，比如多次调用PushMessageById肯定依次到达
    如果在某个前端接口回调前调用的PushMessageById，肯定在前端接口回调前到达
    
    
## Handler和Remote
    处理来自客户端的消息，我们称之为Handler
    后端服务器也可以处理来自客户端的消息，消息会自动转发到对应的后端服务器
    
    处理来自服务器内部服务之间通讯，我们称之为Remote    

    我们定义handler其实主要是通过 服务类型名.handler就是前端协议处理
    比如: 下面
    nodeapi.Registry.AddCollection("chat.handler").
        Register(&Handler{}, apientry.WithName("chat"), apientry.WithNameFunc(strings.ToLower))

    nodeapi.Registry.AddCollection("chat.remote").
		Register(&Entry{}, apientry.WithName("chatremote"), apientry.WithNameFunc(strings.ToLower))
    就分别注册了前后端协议
    
    代码会自动分析方法是否匹配
    需要满足条件，参看 DefaultFormater.IsValidMethod

    Remote的调用方式比如:
    app.Request(sender, "chat.chatremote.entry", param, msg, cb)
    自动根据路由规则来找到目标服务，并发起请求

### Handler消息处理context使用原则
    一般来说消息处理里可能写很多异步代码
    传入的context主要目标也就是传入session    
    在前端服务器，session是FrontSession，在后端服务器，则是BackSession
    Frontsession一般可能在后续被关闭，所以在需要处理Frontsession的场合，可以判断下frontsession.IsClosed
    Frontsession生命周期由外部连接控制   
    BackSession如果想知道通过setValue设置的值，需要先调用QuerySession来获取

## 服务发现模块
    服务发现模块，每个节点会向etcd注册自己和拥有的服务
    做服务转发时，就可以通过路由算法那决定具体向某个服务类型的那个具体服务转发请求

## 路由模块
    路由模块简单来说就是如何在所有的同类服务中选择合适的发送消息
    请参看node/route
    用户可以定义某个服务类型的route函数，route过程通过[serviceType, routeParam]来在
    当前活跃服务中选取目标服务
    定义路由
    func initRoutes() {
        rs := route.GetRouteService()
        rs.Register("chat", func(serviceType string, param route.IRouteParam) string {
            chatid := param.Get("chatid", "").(string)
            return chatid
        })
    }

    参看
    node/appp/serviceutils.go
        Request
    服务器之间的请求请使用此接口
    
    前端服务的FrontSession对象自动作为routeParam用于路由
    参看  func (f *ForwarderComponent) Forward(fs *cs.FrontSession, msg *msgs.ClientMsg)

## FrontSession & BackSession
    前后端协议处理分别有对应的的session
    session可以通过
    fs.Bind(uid)
	fs.Set("chatid", chatItem.Name)
	fs.PushSession()
    来在session上保存数据

    后端服务，必须QuerySession才能拿到当前session数据

## 客户端协议

## NodeService & Component
    最佳实践
    逻辑可以使用Component的方式聚合在NodeService上
    逻辑也可以使用本地Actor/Service的方式提供非阻塞服务

## Master
    _tools/master
    master服务器提供了一个工具，用户可以用telnet连接master来查看服务器组状态，驱动node节点退休

## Hotreplace
    服务器提供退休机制来进行热更，整个热更过程为
    . 退休 node1
    . 等待退休完毕后，要求node1 进程关闭
    . 重启node1

    通过master向节点发送retire指令来启动退休
    节点上的所有service必须都能退休，该节点才可以退休
    如果service要支持退休，请使用SetCtrlCmdListener来处理命令
    处理queryretire命令，返回"ok"
    处理retire，等待retire结束后，向nodeadmin发送

# _examples
## chat
	. 聊天室和roll点功能

# _benchmarks
## rpc
    . 测试rpc基础性能



# 其他
    使用etcd服务器来实现服务发现，请启动etcd
    https://github.com/etcd-io/etcd/releases
    https://github.com/etcd-io/etcd/releases/download/v3.4.19/etcd-v3.4.19-windows-amd64.zip
## chat
### 功能
    可以通过浏览器启动chat-client/index.html来启动客户端
    客户端启动后，登录服务器，随机分配到某个聊天服务器
    并加入房间，房间最多8人，可以聊天
    聊天输入/roll可以掷骰子
    玩家进入服务器总是进入最前面的房间
    房间不会删除
### 配置与启动
    data/config/nodes.yaml
    配置了多个节点
    缺省所有服务都启动在all-1节点上
    也可以启动多个节点node-1-4
    使用id来启动不同的节点类型
    ./chat.exe -id=node-1

## 客户端
    chat-client
    网页客户端 index.html拖拽到浏览器执行

    chat-client-go
    go 客户端，可以接收消息

# refs
https://github.com/NetEase/pomelo
https://github.com/NetEase/pomelo/wiki/Home-in-Chinese
https://github.com/asynkron/protoactor-go/actor
https://github.com/topfreegames/pitaya
https://github.com/lonng/nano