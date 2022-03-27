package components

import "github.com/go-gl/mathgl/mgl64"

type ThirdPersonControllerComponent struct {
	CameraID int

	Controlled bool
	Grounded   bool

	Velocity           mgl64.Vec3
	BaseVelocity       mgl64.Vec3
	ControllerVelocity mgl64.Vec3
	MovementSpeed      float64
}

func (c *ThirdPersonControllerComponent) AddToComponentContainer(container *ComponentContainer) {
	container.ThirdPersonControllerComponent = c
}
