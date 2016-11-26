package render

import (
	"fmt"
	"math"

	"github.com/kkevinchou/ant/lib/math/vector"
)

const (
	cameraRotationXMax float64 = 80
)

type Camera struct {
	Position vector.Vector3
	View     vector.Vector
	speed    float64
}

func NewCamera(position vector.Vector3, view vector.Vector, speed float64) *Camera {
	return &Camera{
		Position: position,
		View:     view,
		speed:    speed,
	}
}

func (c *Camera) ChangeView(v vector.Vector) {
	c.View.X += float64(v.Y) * sensitivity
	c.View.Y += float64(v.X) * sensitivity

	if c.View.X < -cameraRotationXMax {
		c.View.X = -cameraRotationXMax
	}
}

func (c *Camera) Move(v vector.Vector3) {
	forwardVector := c.backward()
	forwardVector = forwardVector.Scale(-v.Z)

	rightVector := c.right()
	rightVector = rightVector.Scale(-v.X)

	c.Position = c.Position.Add(forwardVector).Add(rightVector).Add(vector.Vector3{X: 0, Y: v.Y, Z: 0})
}

func toRadians(degrees float64) float64 {
	return degrees / 180 * math.Pi
}

func (c *Camera) backward() vector.Vector3 {
	xRadianAngle := -toRadians(c.View.X)
	if xRadianAngle < 0 {
		xRadianAngle += 2 * math.Pi
	}
	yRadianAngle := -(toRadians(c.View.Y) - (math.Pi / 2))
	if yRadianAngle < 0 {
		yRadianAngle += 2 * math.Pi
	}

	x := math.Cos(yRadianAngle) * math.Cos(xRadianAngle)
	y := math.Sin(xRadianAngle)
	z := -math.Sin(yRadianAngle) * math.Cos(xRadianAngle)

	return vector.Vector3{X: x, Y: y, Z: z}
}

func (c *Camera) right() vector.Vector3 {
	xRadianAngle := -toRadians(c.View.X)
	if xRadianAngle < 0 {
		xRadianAngle += 2 * math.Pi
	}
	yRadianAngle := -(toRadians(c.View.Y) - (math.Pi / 2))
	if yRadianAngle < 0 {
		yRadianAngle += 2 * math.Pi
	}

	x, y, z := math.Cos(yRadianAngle), math.Sin(xRadianAngle), -math.Sin(yRadianAngle)

	v1 := vector.Vector3{X: x, Y: math.Abs(y), Z: z}
	v2 := vector.Vector3{X: x, Y: 0, Z: z}
	v3 := v1.Cross(v2)

	if v3.X == 0 && v3.Y == 0 && v3.Z == 0 {
		v3 = vector.Vector3{X: v2.Z, Y: 0, Z: -v2.X}
	}

	fmt.Println(v1, v2, v3)

	return v3.Normalize()
}
