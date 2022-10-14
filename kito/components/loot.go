package components

type LootComponent struct {
}

func (c *LootComponent) AddToComponentContainer(container *ComponentContainer) {
	container.LootComponent = c
}

func (c *LootComponent) ComponentFlag() int {
	return ComponentFlagLoot
}
