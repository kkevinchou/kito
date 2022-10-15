package components

import "github.com/kkevinchou/kito/kito/mechanics/items"

type InventoryComponent struct {
	Items []items.Item
}

func (c *InventoryComponent) AddToComponentContainer(container *ComponentContainer) {
	container.InventoryComponent = c
}

func (c *InventoryComponent) ComponentFlag() int {
	return ComponentFlagInventory
}
