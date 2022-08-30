package components

type LootComponent struct {
	Value float64
}

func (c *LootComponent) AddToComponentContainer(container *ComponentContainer) {
	container.LootComponent = c
}

func (c *LootComponent) ComponentFlag() int {
	return ComponentFlagLoot
}
