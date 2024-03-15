package scenem

type SceneObj struct {
	CfgId     int32
	LineId    int32
	ServiceId string
	SceneId   uint64

	Token     int
	TimeStart int64

	// 玩家数量
	PlayerNum int
}

type SceneServiceStat struct {
	ActiveSceneNum int
	CPURate        float32 // 0-1 0代表空闲, 1代表满了
	// 是否提供服务
	Working           bool
	LastActiveTime    int64
	ActiveFailedTimes int
}

// GetBusyWeight 获取本线的繁忙程度(0-1)
func (s *SceneServiceStat) GetBusyWeight() float32 {
	// todo: 后续改成根据cpuRate判定
	if !s.Working {
		return 1
	}
	weight := 0.8*s.CPURate + 0.2*(float32(s.ActiveSceneNum)/1000)
	if weight > 1 {
		weight = 1
	}
	return weight
}
