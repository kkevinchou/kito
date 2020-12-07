package components

import (
	"math"

	"github.com/kkevinchou/kito/lib/math/vector"
)

type FreeViewComponent struct {
	view vector.Vector
}

func (c *FreeViewComponent) View() vector.Vector {
	return c.view
}

func (c *FreeViewComponent) SetView(view vector.Vector) {
	c.view = view
}

func (c *FreeViewComponent) UpdateView(delta vector.Vector) {
	c.view.X += delta.Y * ySensitivity
	c.view.Y += delta.X * xSensitivity

	if c.view.X < -cameraRotationXMax {
		c.view.X = -cameraRotationXMax
	} else if c.view.X > cameraRotationXMax {
		c.view.X = cameraRotationXMax
	}
}

func (c *FreeViewComponent) Forward() vector.Vector3 {
	xRadianAngle := -toRadians(c.view.X)
	if xRadianAngle < 0 {
		xRadianAngle += 2 * math.Pi
	}
	yRadianAngle := -(toRadians(c.view.Y) - (math.Pi / 2))
	if yRadianAngle < 0 {
		yRadianAngle += 2 * math.Pi
	}

	x := math.Cos(yRadianAngle) * math.Cos(xRadianAngle)
	y := math.Sin(xRadianAngle)
	z := -math.Sin(yRadianAngle) * math.Cos(xRadianAngle)

	return vector.Vector3{X: x, Y: y, Z: z}.Scale(-1)
}

func (c *FreeViewComponent) Right() vector.Vector3 {
	xRadianAngle := -toRadians(c.view.X)
	if xRadianAngle < 0 {
		xRadianAngle += 2 * math.Pi
	}
	yRadianAngle := -(toRadians(c.view.Y) - (math.Pi / 2))
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

	return v3.Normalize()
}
