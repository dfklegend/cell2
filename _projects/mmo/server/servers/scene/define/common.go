package define

import (
	"mmo/libs/simonwittber/go-vector"
)

type UnitId int32

// UnitType entity单元类型
type UnitType int32
type Pos = vector.Vector3

const (
	InvalidUnitId UnitId  = 0
	MaxWidth      float32 = 30 // 场景缺省尺寸
)
