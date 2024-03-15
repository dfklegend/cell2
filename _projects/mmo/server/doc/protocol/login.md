# 说明
    登录相关协议

# QueryGate
    Name: query gate server
    Desc: 请求分配一个gate服务器建立连接
    Route: querygate
    Forward: client -> gate
    Args: cproto.QueryGateReq
    Ret: cproto.QueryGateAck

# db认证
    Name: 认证协议
    Desc: 数据库认证，获得角色id
    Route: db.dbremote.auth
    Forward: gate -> db
    Args: DBAuth
    Ret: DBAuthAck

