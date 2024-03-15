package sceneservice

import (
	"mmo/common/entity"
	"mmo/servers/scene/define"
)

// OnAddEntity 可添加一些快捷方式
// TODO: enter and leave
func (s *Scene) OnAddEntity(e entity.IEntity) {
	s.logic.OnAddEntity(e)
}

func (s *Scene) OnDestroyEntity(e entity.IEntity) {
	s.logic.OnDestroyEntity(e)
}

func (s *Scene) OnCameraEnter(camera define.ICamera) {
	// . 遍历可见对象
	// . 收集可见对象的快照
	// . 发送给camera

	s.logic.OnPreCameraEnter(camera)
	s.pushAllSnapshots(camera)
	s.logic.OnPostCameraEnter(camera)
}

func (s *Scene) pushAllSnapshots(camera define.ICamera) {
	s.world.Filter(func(e entity.IEntity) {
		s.logic.PushSnapshot(camera, e)
	})
}

func (s *Scene) AddCamera(id entity.EntityID, camera define.ICamera) {
	s.cameras[id] = camera
}

func (s *Scene) RemoveCamera(id entity.EntityID) {
	delete(s.cameras, id)
}

func (s *Scene) PushViewSnapshot(e entity.IEntity) {
	for _, v := range s.cameras {
		s.logic.PushSnapshot(v, e)
	}
}

func (s *Scene) GetEntity(id entity.EntityID) entity.IEntity {
	return s.world.GetEntity(id)
}
