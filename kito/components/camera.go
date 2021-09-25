package components

// Mostly work as a type flag atm
type CameraComponent struct {
}

func (c *CameraComponent) AddToComponentContainer(container *ComponentContainer) {
	container.CameraComponent = c
}
