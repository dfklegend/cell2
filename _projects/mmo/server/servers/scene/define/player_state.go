package define

// ScenePlayerState 玩家状态
type ScenePlayerState int

const (
	Init          ScenePlayerState = iota
	LoadingInfo                    // 载入数据
	PreNormal                      //准备状态，向logic再确认一次成功即转化为Normal，确认失败则是WaitMgrDelete
	Normal                         // 正常 (离线单独标志)
	Logouting                      // 登出中
	WaitMgrDelete                  // 等待删除
)
