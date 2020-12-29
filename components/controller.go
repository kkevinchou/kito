package components

type ControllerComponent struct {
	controlled bool
}

func NewControllerComponent() *ControllerComponent {
	component := ControllerComponent{}
	return &component
}

func (c *ControllerComponent) SetControlled(controlled bool) {
	c.controlled = controlled
}

func (c *ControllerComponent) Controlled() bool {
	return c.controlled
}

func (c *ControllerComponent) AddToComponentContainer(container *ComponentContainer) {
	container.ControllerComponent = c
}
