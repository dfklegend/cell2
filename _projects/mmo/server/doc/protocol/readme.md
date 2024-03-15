# 说明
    将协议整理出来，便于维护
    可以根据逻辑整理多个md文件
    协议规范
    
    Name: 名字
    Desc: 协议用途
    Route: 路由
    Forward: 发送方向
    Args: 参数类型
    Ret: 返回值

    举例
    Name: 认证协议
    Desc: 数据库认证，获得角色id
    Route: db.dbremote.auth
    Forward: gate -> db
    Args: DBAuth
    Ret: DBAuthAck

