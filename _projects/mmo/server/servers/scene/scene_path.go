package sceneservice

import (
	"github.com/dfklegend/astar"

	"mmo/servers/scene/define"
	"mmo/utils"
)

func (s *Scene) ToGridX(x float32) int {
	return utils.ToGridX(x, s.width)
}

func (s *Scene) GridToX(x int) float32 {
	return utils.GridToX(x, s.width)
}

func (s *Scene) ToGridZ(z float32) int {
	return utils.ToGridX(z, s.height)
}

func (s *Scene) GridToZ(z int) float32 {
	return utils.GridToX(z, s.height)
}

func (s *Scene) FindPath(src, tar define.Pos) *define.Path {
	srcX := s.ToGridX(src.X)
	srcZ := s.ToGridZ(src.Z)

	tarX := s.ToGridX(tar.X)
	tarZ := s.ToGridZ(tar.Z)

	s.pathContext.AddTempSkipPoint(tarX, tarZ)
	//s.ns.GetLogger().Infof("findPath (%v, %v) -> (%v, %v)", src.X, src.Z, tar.X, tar.Z)
	nodes, err := s.finder.FindPathEx(s.pathContext, astar.Node{X: srcX, Y: srcZ},
		astar.Node{X: tarX, Y: tarZ}, 99)
	s.pathContext.ClearTempSkipPoints()
	if err != nil {
		s.ns.GetLogger().Infof("%v ret nil path (%v, %v) -> (%v, %v)",
			s.sceneId,
			srcX, srcZ, tarX, tarZ)
		return nil
	}

	if len(nodes) < 2 {
		return nil
	}

	path := define.NewPath()
	// revert it
	// skip start pos
	for i := len(nodes) - 1; i >= 0; i-- {
		v := nodes[i]
		pos := define.Pos{
			X: s.GridToX(v.X),
			Z: s.GridToZ(v.Y),
		}
		path.Points = append(path.Points, pos)
	}

	return path
}

func (s *Scene) gridIsValid(gridX, gridZ int) bool {
	if gridX < 0 || gridX >= s.width {
		return false
	}
	if gridZ < 0 || gridZ >= s.height {
		return false
	}
	return true
}

func (s *Scene) IsValidPos(tar define.Pos) bool {
	// 避免大于场景
	gridX := s.ToGridX(tar.X)
	gridZ := s.ToGridZ(tar.Z)
	if !s.gridIsValid(gridX, gridZ) {
		return false
	}
	return true
}

func (s *Scene) IsInBlock(tar define.Pos) bool {
	gridX := s.ToGridX(tar.X)
	gridZ := s.ToGridZ(tar.Z)
	if !s.gridIsValid(gridX, gridZ) {
		return false
	}
	if s.graph.IsInBlock(gridX, gridZ) {
		return true
	}
	return false
}
