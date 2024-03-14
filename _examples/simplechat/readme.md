# 实现聊天测试

    当前客户端没有实现，直接使用节点作为客户端

# 简单功能列表

    . 加入聊天
    . 修改昵称
    . 切换聊天室
    . roll点流程

# 实现

	使用ClientService来模拟客户端
	连接到gate请求分配一个connector
	发送登录到connector
	connector会请求chatMgr来分配一个房间进入
	成功进入后，玩家也可以请求 进入另外一个房间
	玩家进入房间后，保持ping来keepalive


	

	服务类型
		client
			客户端
		gate
			connector分配
		connector
			连接器
		login
			登录
		chatm
			聊天管理
		chat
			聊天

	
    ClientService
        申请分配聊天室
        请求进入聊天室
        聊天
        定时向聊天室发送ping
    
    GateService
        分配Connector
        ReqConnector
    
    ConnectorService
        连接服务
        后期监听客户端连接
        Login
        AckLogin
    
    LoginService
        分配一个唯一id
        LSLogin
        
        LSAckLogin
            Id int32

    ChatMgr
        负责分配chatservice

        CMReqLogin
            请求登录到聊天系统

        CMAckLogin
            分配的chatservice
            分配的房间id
            token

    chatservice
        负责管理实际的聊天房间

		CSCreateRoom
		
		CSAckCreateRoom

        CSReqLogin
            请求进入
            已有房间或者新建房间

        CSAckLogin

        CSSay
            
        CSReqChangeName

        CSPing

# 流程

## 登录

	client		connector		chatm		chat
				client->login

								connector->CMReqLogin
											chatm->CSCreateRoom
								<-
				<-CMAckLogin
	<-AckLogin

## 房间进入

	client可以通过登录获取房间令牌
	获得令牌后，直接向房间服务登录(注册自己的PID，便于消息广播)
	房间成员有变化后，都会向chatmgr同步状态，便于人数控制

## Ping

	进入房间后，客户端定期向房间服务发送ping
	房间发现某个成员长时间未ping，将移除玩家