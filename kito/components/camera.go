package components

// Mostly work as a type flag atm
type CameraComponent struct {
}

func (c *CameraComponent) AddToComponentContainer(container *ComponentContainer) {
	container.CameraComponent = c
}

func (c *CameraComponent) ComponentFlag() int {
	return ComponentFlagCamera
}
