//Package vector provides Vector2, Vector3 and Vector4 types with the
//usual family of vector functions,
package vector

import "math"

type Vector2 struct {
	X float32
	Y float32
}

func (a Vector2) Magnitude() float32 {
	return float32(math.Sqrt(float64(a.X*a.X + a.Y*a.Y)))
}

func (a Vector2) SqrMagnitude() float32 {
	return a.X*a.X + a.Y*a.Y
}

func (a Vector2) Dot(b Vector2) float32 {
	return a.X*b.X + a.Y*b.Y
}

func (a Vector2) Sub(b Vector2) Vector2 {
	return Vector2{a.X - b.X, a.Y - b.Y}
}

func (a Vector2) Add(b Vector2) Vector2 {
	return Vector2{a.X + b.X, a.Y + b.Y}
}

func (a Vector2) Mul(b float32) Vector2 {
	return Vector2{a.X * b, a.Y * b}
}

func (a Vector2) Div(b float32) Vector2 {
	return Vector2{a.X / b, a.Y / b}
}

func (a Vector2) Scale(b Vector2) Vector2 {
	return Vector2{a.X * b.X, a.Y * b.Y}
}

func (a Vector2) Project(onNormal Vector2) Vector2 {
	n := onNormal.Dot(onNormal)
	if n < 1.e-5 {
		return Vector2{0, 0}
	}
	return onNormal.Mul(a.Dot(onNormal)).Div(n)
}

func (a Vector2) Reflect(inNormal Vector2) Vector2 {
	return inNormal.Mul(-2 * inNormal.Dot(a)).Add(a)
}

func (a Vector2) Normalized() Vector2 {
	magnitude := a.Magnitude()
	if magnitude > 1.e-5 {
		return a.Div(magnitude)
	} else {
		return Vector2{0, 0}
	}
}

func (a Vector2) Distance(b Vector2) float32 {
	return a.Sub(b).Magnitude()
}

func (a Vector2) Lerp(b Vector2, d float32) Vector2 {
	return Vector2{a.X + (b.X-a.X)*d, a.Y + (b.Y-a.Y)*d}
}
