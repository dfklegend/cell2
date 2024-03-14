package service

/*
	service包
	service提供独立的routine来处理事务
	基于actor提供服务，将请求转发给对应的接口
	接口处理在service所在routine处理
	调用Request需要在service所在routine调用，可以Post进来

	Request根据route来选择接口进行处理
	Request CB回来的消息，由protobuf打包(附带了类型名)，返回来后由调用Service序列化成目标类型，交给CB
	(依赖于protobuf的类型名机制)

	TODO:
	service应该不要自动重启，crash就关闭
*/
