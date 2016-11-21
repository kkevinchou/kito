package vector

import (
	"fmt"
	"math"
)

type Vector struct {
	X float64
	Y float64
}

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func Zero() Vector {
	return Vector{X: 0, Y: 0}
}

func (v Vector) String() string {
	return fmt.Sprintf("<Vector %.2f, %.2f>", v.X, v.Y)
}

func (v1 Vector) Add(v2 Vector) Vector {
	return Vector{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
}

func (v1 Vector) Sub(v2 Vector) Vector {
	return Vector{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
	}
}

func (v Vector) Normalize() Vector {
	return v.Scale(1.0 / v.Length())
}

func (v Vector) Scale(s float64) Vector {
	return Vector{
		X: v.X * s,
		Y: v.Y * s,
	}
}

func (v Vector) Length() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}

func (v Vector) Clamp(max float64) Vector {
	length := v.Length()

	if length > max {
		return v.Scale(max / length)
	} else {
		return v
	}
}

func (v Vector) Cross(v2 Vector) float64 {
	return (v.X * v2.Y) - (v.Y * v2.X)
}

func (v Vector3) Cross(other Vector3) Vector3 {
	return Vector3{
		X: v.Y*other.Z - other.Y*v.Z,
		Y: -(v.X*other.Z - other.X*v.Z),
		Z: v.X*other.Y - other.X*v.Y,
	}
}

func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

func (v Vector3) Length() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2))
}

func (v Vector3) Normalize() Vector3 {
	return v.Scale(1.0 / v.Length())
}

func (v Vector3) Scale(s float64) Vector3 {
	return Vector3{
		X: v.X * s,
		Y: v.Y * s,
		Z: v.Z * s,
	}
}

func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

func (v Vector3) Clamp(max float64) Vector3 {
	length := v.Length()

	if length > max {
		return v.Scale(max / length)
	} else {
		return v
	}
}
