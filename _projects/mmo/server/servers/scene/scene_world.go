package sceneservice

import (
	"github.com/dfklegend/cell2/node/app"

	"mmo/servers/scene/define"
)

// PushViewMsg 下发视图消息
func (s *Scene) PushViewMsg(pos define.Pos, route string, msg any) {
	// 后续将改成依赖于viewStrategy
	// TODO: 优化，多个camera如果在同一个front，会合并
	// camera组织成利于分发的形式
	for _, v := range s.cameras {
		app.PushMessageById(s.ns, v.GetFrontId(), v.GetNetId(), route, msg)
	}
}
