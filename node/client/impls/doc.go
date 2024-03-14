package impls

/*
如何扩展协议

必须要有 route字段, 用来进行协议路由 必须式A.B.C格式，A代表服务类型，B代表接口类别，C代表函数名
数据字段使用json或者 protobuf来序列化
如果不需要Request/Response方式，可以没有clientReqId字段
*/
