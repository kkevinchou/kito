package components

type ThirdPersonControllerComponent struct {
	Controlled bool
	CameraID   int
}

func (c *ThirdPersonControllerComponent) AddToComponentContainer(container *ComponentContainer) {
	container.ThirdPersonControllerComponent = c
}
