package kito

import (
	"math"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/math/matrix"
	"github.com/kkevinchou/kito/lib/math/vector"
)

const (
	cameraRotationXMax = 80
	cameraSpeedScalar  = 10
	sensitivity        = 0.3
)

type Followable interface {
	Position() vector.Vector3
}

type Camera struct {
	position     vector.Vector3
	view         vector.Vector
	speed        vector.Vector3
	followTarget Followable

	// vector indicating what the camera is being commanded to do.
	// Forward, Backward, Left, Right, Up
	// Encoded as a directional vector.
	// This could potentially be an enum as we only have 5 valid states
	commandHeading vector.Vector3

	controlled bool
}

func NewCamera(position vector.Vector3, view vector.Vector) *Camera {
	return &Camera{
		position: position,
		view:     view,
	}
}

func (c *Camera) ChangeView(v vector.Vector) {
	c.view.X += float64(v.Y) * sensitivity
	c.view.Y += float64(v.X) * sensitivity

	if c.view.X < -cameraRotationXMax {
		c.view.X = -cameraRotationXMax
	}

	if c.view.X > cameraRotationXMax {
		c.view.X = cameraRotationXMax
	}
}

func (c *Camera) SetCommandHeading(v vector.Vector3) {
	c.commandHeading = v

	forwardVector := c.backward()
	forwardVector = forwardVector.Scale(-c.commandHeading.Z)

	rightVector := c.right()
	rightVector = rightVector.Scale(-c.commandHeading.X)
	c.speed = forwardVector.Add(rightVector).Add(vector.Vector3{X: 0, Y: c.commandHeading.Y, Z: 0})
}

func toRadians(degrees float64) float64 {
	return degrees / 180 * math.Pi
}

func (c *Camera) backward() vector.Vector3 {
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

	return vector.Vector3{X: x, Y: y, Z: z}
}

func (c *Camera) right() vector.Vector3 {
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

func (c *Camera) Update(delta time.Duration) {
	if c.followTarget != nil {
		c.position = c.followTarget.Position()
		c.position.Y += 5
		c.position.Z += 5
		return
	}

	if c.speed == vector.Zero3() {
		return
	}
	c.position = c.position.Add(c.speed.Normalize().Scale(cameraSpeedScalar * delta.Seconds()))
}

func (c *Camera) Position() vector.Vector3 {
	return c.position
}

func (c *Camera) View() vector.Vector {
	return c.view
}

func (c *Camera) GetRayDirection(x, y float64) vector.Vector3 {
	// Get the projection matrix
	pMatrixValues := make([]float32, 16)
	gl.GetFloatv(gl.PROJECTION_MATRIX, &pMatrixValues[0])

	mvMatrixValues := make([]float32, 16)
	gl.GetFloatv(gl.MODELVIEW, &mvMatrixValues[0])

	mvMatrix := matrix.Mat4FromValues(mvMatrixValues)
	pMatrix := matrix.Mat4FromValues(pMatrixValues)

	// Convert the screen coordinate to normalised device coordinates
	NDCPoint := mgl32.Vec4{(2.0*float32(x))/800 - 1, 1 - (2.0*float32(y))/600, -1, 1}
	worldPoint := pMatrix.Mul4(mvMatrix).Inv().Mul4x1(NDCPoint)

	// Normalize on W
	worldPoint = mgl32.Vec4{worldPoint[0] / worldPoint[3], worldPoint[1] / worldPoint[3], worldPoint[2] / worldPoint[3], 1}

	// Extract the 3D vector
	worldPointVector := vector.Vector3{X: float64(worldPoint[0]), Y: float64(worldPoint[1]), Z: float64(worldPoint[2])}
	return worldPointVector.Sub(c.Position()).Normalize()
}

func (c *Camera) Follow(entity Followable) {
	c.followTarget = entity
}
