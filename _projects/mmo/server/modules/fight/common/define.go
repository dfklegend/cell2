package common

import (
	"mmo/common/entity"
	"mmo/libs/simonwittber/go-vector"
)

type BufId = string
type SkillId = string
type EquipId = string

type Pos = vector.Vector3
type CharId = entity.EntityID

// 接口定义，避免交叉引用
