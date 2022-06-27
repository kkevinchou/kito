package netsync

import "github.com/go-gl/mathgl/mgl64"

const (
	gravity float64 = 250
)

var (
	accelerationDueToGravity = mgl64.Vec3{0, -gravity, 0}
)
