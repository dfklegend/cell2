# 说明
    有部分需要打开ectd服务器

# test-service
    测试服务之间的Request调用
    client -> addservice1 -> addservice2

# test-cluster
    测试cluster的组织

# test-app
    使用service来组织app

# simplechat
    在没有client协议的情况下，使用client服务来模拟一个聊天实现来检验框架

# chatapi
    使用接口映射来实现simplechat
    (逻辑流程参看simplechat readme.md)
    测试，启动后
    在控制台输入
    /login 1
    /say ffff    

    . 向gate 分配connector
    . connector先向login请求登录(分配uid)，再向chatm登录
    . chatm分配个空房间(在那个chatservice,房间id等)
    . 进入房间(使用token)

# chat2
    完整版的cell2服务应用
    前端handler和后端remote接口使用
    多个chat服务，则随机分配一个服务加入房间
    一个房间8个人，满了创建新房间

# chat2-client-js
    chat2 js的客户端
    可以输入聊天，可以使用/roll 来掷骰子 同房间可以看到消息

# chat2-client-go
    go语言客户端
    不能输入
    只是进入房间，可以收到web客户端来的消息

# test-msg-guaranteed
    测试如果目标服务器暂时不可用，消息是否会丢失
	结果是，消息会在连接上后，继续发送
	(最大限度的会重发，当然内存会有消耗)
	(极端情况下，丢失，接收者收到后，进程crash了，会丢失)