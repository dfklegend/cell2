package sceneplayer

/*
	场景玩家对象
	由logic创建
	异常处理
	创建
		player成功进入世界后，二次向logic确认一次，确认成功则正式进入normal状态，否则认为是异常，等待删除
	keepalive
		player定时向logic确认 keepalive
		如果失败
			// 启动异常处理
*/
