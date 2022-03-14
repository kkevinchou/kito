package components

import "github.com/go-gl/mathgl/mgl64"

type ThirdPersonControllerComponent struct {
	Controlled     bool
	CameraID       int
	MovementVector mgl64.Vec3
	Grounded       bool

	Velocity     mgl64.Vec3
	BaseVelocity mgl64.Vec3
}

func (c *ThirdPersonControllerComponent) AddToComponentContainer(container *ComponentContainer) {
	container.ThirdPersonControllerComponent = c
}
