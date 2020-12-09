package components

type PlayerControllerComponent struct {
	controlled bool
}

func NewPlayerControllerComponent() *PlayerControllerComponent {
	component := PlayerControllerComponent{}
	return &component
}

func (c *PlayerControllerComponent) SetControlled(controlled bool) {
	c.controlled = controlled
}

func (c *PlayerControllerComponent) Controlled() bool {
	return c.controlled
}
