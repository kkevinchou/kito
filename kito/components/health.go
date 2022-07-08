package components

type HealthComponent struct {
	Value float64
}

func (c *HealthComponent) AddToComponentContainer(container *ComponentContainer) {
	container.HealthComponent = c
}

func (c *HealthComponent) ComponentFlag() int {
	return ComponentFlagHealth
}
