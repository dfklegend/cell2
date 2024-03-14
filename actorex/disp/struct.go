package disp

//	disp包
//	proto.actor的调度器实现
//	提供可控的携程环境

type ChanTask chan func()
