package define

type Path struct {
	Points []Pos
}

func NewPath() *Path {
	return &Path{
		Points: make([]Pos, 0, 10),
	}
}
