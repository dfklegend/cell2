package scenem

import (
	"github.com/dfklegend/cell2/node/app"
	l "github.com/dfklegend/cell2/utils/logger"

	mymsg "mmo/messages"
)

func (m *SceneServiceMgr) SpawnScene(cfgId int32) bool {
	info := m.AllocScene(cfgId)
	if info == nil {
		return false
	}

	app.Request(m.ns, "scene.remote.allocscene", info.ServiceId, &mymsg.SAllocScene{
		CfgId:   cfgId,
		SceneId: info.SceneId,
		Token:   int32(info.Token),
	}, func(err error, raw interface{}) {
		if err != nil {
			l.L.Errorf("scene.remote.allocscene error: %v", err)
			return
		}

		m.OnSceneCreateSucc(info)
	})
	return true
}
