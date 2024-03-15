package components

import (
	"mmo/servers/scene/blockgraph"
	"mmo/servers/scene/define"
)

type Transform struct {
	*BaseSceneComponent

	pos  define.Pos
	dirt bool

	block *blockgraph.Holder
}

func NewTransform() *Transform {
	return &Transform{
		BaseSceneComponent: NewBaseSceneComponent(),
	}
}

func (t *Transform) OnPrepare() {
	t.BaseSceneComponent.OnPrepare()
	t.block = blockgraph.NewHolder(t.scene.GetGraph())
}

func (t *Transform) OnStart() {
}

func (t *Transform) Update() {
}

func (t *Transform) LateUpdate() {
	t.checkUpdatePos()
}

func (t *Transform) OnDestroy() {
	t.block.Clear()
}

func (t *Transform) GetPos() define.Pos {
	return t.pos
}

func (t *Transform) SetPos(pos define.Pos) {
	if t.pos.Distance(pos) < 0.001 {
		return
	}
	t.pos = pos
	t.updateBlockPos()

	t.dirt = true
}

func (t *Transform) checkUpdatePos() {
	if !t.dirt {
		return
	}
	t.dirt = false

	t.updateSpacePos()
	t.updateBlockPos()
}

func (t *Transform) updateSpacePos() {
	space := t.GetScene().GetSpace()
	space.UpdateEntityPos(t.GetOwner().GetId(), t.pos)
}

func (t *Transform) updateBlockPos() {
	if t.block == nil {
		return
	}
	t.block.Update(t.scene.ToGridX(t.pos.X), t.scene.ToGridZ(t.pos.Z))
}

func (t *Transform) ClearBlockPos() {
	t.block.Clear()
}

func (t *Transform) ResetBlockPos() {
	t.updateBlockPos()
}
