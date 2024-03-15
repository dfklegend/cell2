package bridge

// player处理时，有时候需要明确的调用system的函数，这里定义接口
// 通过接口，来避免交叉引用
// system通过ILogicPlayer来访问player
