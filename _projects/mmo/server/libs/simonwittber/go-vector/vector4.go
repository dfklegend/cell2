package vector

import "math"

type Vector4 struct {
	W float32
	X float32
	Y float32
	Z float32
}

func (a Vector4) Magnitude() float32 {
	return float32(math.Sqrt(float64(a.W*a.W + a.X*a.X + a.Y*a.Y + a.Z*a.Z)))
}

func (a Vector4) SqrMagnitude() float32 {
	return a.W*a.W + a.X*a.X + a.Y*a.Y + a.Z*a.Z
}

func (a Vector4) Dot(b Vector4) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W
}

func (a Vector4) Sub(b Vector4) Vector4 {
	return Vector4{a.W - b.W, a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (a Vector4) Add(b Vector4) Vector4 {
	return Vector4{a.W + b.W, a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a Vector4) Mul(b float32) Vector4 {
	return Vector4{a.W * b, a.X * b, a.Y * b, a.Z * b}
}

func (a Vector4) Div(b float32) Vector4 {
	return Vector4{a.W / b, a.X / b, a.Y / b, a.Z / b}
}

func (a Vector4) Scale(b Vector4) Vector4 {
	return Vector4{a.W * b.W, a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

func (a Vector4) Project(b Vector4) Vector4 {
	return b.Mul(a.Dot(b) / b.Dot(b))
}

func (a Vector4) Normalized() Vector4 {
	magnitude := a.Magnitude()
	if magnitude > 1.e-5 {
		return a.Div(magnitude)
	} else {
		return Vector4{0, 0, 0, 0}
	}
}

func (a Vector4) Distance(b Vector4) float32 {
	return a.Sub(b).Magnitude()
}

func (a Vector4) Lerp(b Vector4, d float32) Vector4 {
	return Vector4{a.W + (b.W-a.W)*d, a.X + (b.X-a.X)*d, a.Y + (b.Y-a.Y)*d, a.Z + (b.Z-a.Z)*d}
}
