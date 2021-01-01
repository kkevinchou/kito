package components

type ControllerComponent struct {
	Controlled  bool
	IsCharacter bool
}

func (c *ControllerComponent) AddToComponentContainer(container *ComponentContainer) {
	container.ControllerComponent = c
}
