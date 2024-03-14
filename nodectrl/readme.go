package nodectrl

/*
	管理节点
		主要为热替换服务
		Cell2内，热替换以节点为单位，节点内所有service必须都支持热替换才允许节点热替换

	master -> nodeadmin
	stat
		查看状态
	retire
		退休本节点

		退休命令
		设置状态为 retiring
		向每个service发送 retire命令
		每个service退休完毕，则通知admin retired
		则node为退休状态(retired)

	exit
		retired状态才处理exit消息，退出自身


	service支持的命令
	nodeadmin -> service
	queryretire
		查询是否支持退休
	retire
		开始退休

	nodeadmin <- service
	(servicecmd)
	retired
		已退休完毕

*/
