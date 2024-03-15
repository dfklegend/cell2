PACKAGE DOCUMENTATION

package vector
    import "."

    Package vector provides Vector2, Vector3 and Vector4 types with the
    usual family of vector functions,

TYPES

type Vector2 struct {
    X float32
    Y float32
}

func (a Vector2) Add(b Vector2) Vector2

func (a Vector2) Distance(b Vector2) float32

func (a Vector2) Div(b float32) Vector2

func (a Vector2) Dot(rhs Vector2) float32

func (a Vector2) Lerp(b Vector2, d float32) Vector2

func (a Vector2) Magnitude() float32

func (a Vector2) Mul(b float32) Vector2

func (a Vector2) Normalized() Vector2

func (a Vector2) Project(onNormal Vector2) Vector2

func (a Vector2) Reflect(inNormal Vector2) Vector2

func (a Vector2) Scale(b Vector2) Vector2

func (a Vector2) SqrMagnitude() float32

func (a Vector2) Sub(b Vector2) Vector2

type Vector3 struct {
    X float32
    Y float32
    Z float32
}

func (a Vector3) Add(b Vector3) Vector3

func (a Vector3) Cross(rhs Vector3) Vector3

func (a Vector3) Distance(b Vector3) float32

func (a Vector3) Div(b float32) Vector3

func (a Vector3) Dot(rhs Vector3) float32

func (a Vector3) Lerp(b Vector3, d float32) Vector3

func (a Vector3) Magnitude() float32

func (a Vector3) Mul(b float32) Vector3

func (a Vector3) Normalized() Vector3

func (a Vector3) Project(onNormal Vector3) Vector3

func (a Vector3) ProjectOnPlane(planeNormal Vector3) Vector3

func (a Vector3) Reflect(inNormal Vector3) Vector3

func (a Vector3) Scale(b Vector3) Vector3

func (a Vector3) SqrMagnitude() float32

func (a Vector3) Sub(b Vector3) Vector3

type Vector4 struct {
    W float32
    X float32
    Y float32
    Z float32
}

func (a Vector4) Add(b Vector4) Vector4

func (a Vector4) Distance(b Vector4) float32

func (a Vector4) Div(b float32) Vector4

func (a Vector4) Dot(b Vector4) float32

func (a Vector4) Lerp(b Vector4, d float32) Vector4

func (a Vector4) Magnitude() float32

func (a Vector4) Mul(b float32) Vector4

func (a Vector4) Normalized() Vector4

func (a Vector4) Project(b Vector4) Vector4

func (a Vector4) Scale(b Vector4) Vector4

func (a Vector4) SqrMagnitude() float32

func (a Vector4) Sub(b Vector4) Vector4


