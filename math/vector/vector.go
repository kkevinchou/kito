package vector

import (
	"fmt"
	"math"
)

type Vector struct {
	X float64
	Y float64
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
