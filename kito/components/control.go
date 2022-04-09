package components

// mostly just holds the id of the player that controls this thing
type ControlComponent struct {
	PlayerID int
}

func (c *ControlComponent) AddToComponentContainer(container *ComponentContainer) {
	container.ControlComponent = c
}

func (c *ControlComponent) ComponentFlag() int {
	return ComponentFlagControl
}
