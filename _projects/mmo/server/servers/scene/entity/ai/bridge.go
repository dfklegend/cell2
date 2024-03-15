package ai

import (
	"mmo/common/entity"
)

// aictrl需要components的细节
// aictrl又在components内实例化，解决交叉

type IAICtrl interface {
	InitCtrl(owner entity.IEntity)
	OnPrepare()
	OnStart()
	Update()
	SetActive(b bool)
}
