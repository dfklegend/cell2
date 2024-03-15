package blockgraph

type IBlockGraph interface {
	IncBlock(x, y int)
	DecBlock(x, y int)
	IsInBlock(x, y int) bool
}

// Graph 阻挡图
type Graph struct {
	blocks []int8
	width  int
	height int
}

func NewGraph(width, height int) *Graph {
	return &Graph{
		width:  width,
		height: height,
		blocks: make([]int8, width*height),
	}
}

func (g *Graph) offBlock(x, y, off int) {
	if !g.isValidPos(x, y) {
		return
	}
	index := y*g.width + x

	g.blocks[index] += int8(off)
}

func (g *Graph) IncBlock(x, y int) {
	g.offBlock(x, y, 1)
}

func (g *Graph) DecBlock(x, y int) {
	g.offBlock(x, y, -1)
}

func (g *Graph) IsInBlock(x, y int) bool {
	if !g.isValidPos(x, y) {
		return true
	}
	index := y*g.width + x

	return g.blocks[index] > 0
}

func (g *Graph) isValidPos(x, y int) bool {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return false
	}
	return true
}
