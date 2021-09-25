package components

// Mostly work as a type flag atm
type NetworkComponent struct {
}

func (c *NetworkComponent) AddToComponentContainer(container *ComponentContainer) {
	container.NetworkComponent = c
}
