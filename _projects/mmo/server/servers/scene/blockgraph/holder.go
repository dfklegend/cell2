package blockgraph

type Holder struct {
	x      int
	y      int
	setted bool

	graph *Graph
}

func NewHolder(graph *Graph) *Holder {
	return &Holder{
		setted: false,
		graph:  graph,
	}
}

func (h *Holder) Update(x, y int) {
	if h.setted {
		if x == h.x && y == h.y {
			return
		}
		h.graph.DecBlock(h.x, h.y)
	} else {
		h.setted = true
	}

	h.graph.IncBlock(x, y)
	h.x = x
	h.y = y
}

func (h *Holder) Clear() {
	if h.setted {
		h.graph.DecBlock(h.x, h.y)
		h.x = 0
		h.y = 0
		h.setted = false
	}
}
