package render

import (
	"math"

	"github.com/kkevinchou/ant/lib/math/vector"
)

func (r *RenderSystem) CameraView(x, y int) {
	cameraRotationY += float64(x) * sensitivity
	cameraRotationX += float64(y) * sensitivity

	if cameraRotationX < -cameraRotationXMax {
		cameraRotationX = -cameraRotationXMax
	}

	if cameraRotationX > cameraRotationXMax {
		cameraRotationX = cameraRotationXMax
	}
}

func (r *RenderSystem) MoveCamera(v vector.Vector3) {
	forwardX, forwardY, forwardZ := forward()
	// Moving backwards
	forwardX *= -v.Z
	forwardY *= -v.Z
	forwardZ *= -v.Z

	rightX, rightY, rightZ := right()
	rightX *= -v.X
	rightY *= -v.X
	rightZ *= -v.X

	cameraX += forwardX + rightX
	cameraY += forwardY + rightY + v.Y
	cameraZ += forwardZ + rightZ
}

func toRadians(degrees float64) float64 {
	return degrees / 180 * math.Pi
}

func forward() (float64, float64, float64) {
	xRadianAngle := -toRadians(cameraRotationX)
	if xRadianAngle < 0 {
		xRadianAngle += 2 * math.Pi
	}
	yRadianAngle := -(toRadians(cameraRotationY) - (math.Pi / 2))
	if yRadianAngle < 0 {
		yRadianAngle += 2 * math.Pi
	}

	x := math.Cos(yRadianAngle) * math.Cos(xRadianAngle)
	y := math.Sin(xRadianAngle)
	z := -math.Sin(yRadianAngle) * math.Cos(xRadianAngle)

	return x, y, z
}

func right() (float64, float64, float64) {
	xRadianAngle := -toRadians(cameraRotationX)
	if xRadianAngle < 0 {
		xRadianAngle += 2 * math.Pi
	}
	yRadianAngle := -(toRadians(cameraRotationY) - (math.Pi / 2))
	if yRadianAngle < 0 {
		yRadianAngle += 2 * math.Pi
	}

	x, y, z := math.Cos(yRadianAngle), math.Sin(xRadianAngle), -math.Sin(yRadianAngle)

	v1 := vector.Vector3{x, math.Abs(y), z}
	v2 := vector.Vector3{x, 0, z}
	v3 := v1.Cross(v2)

	if v3.X == 0 && v3.Y == 0 && v3.Z == 0 {
		v3 = vector.Vector3{v2.Z, 0, -v2.X}
	}

	v3 = v3.Normalize()

	return v3.X, v3.Y, v3.Z
}
