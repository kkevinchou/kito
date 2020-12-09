package components

import (
	"fmt"

	"github.com/kkevinchou/kito/types"
)

type InventoryComponent struct {
	items map[int]types.Item
}

func NewInventoryComponent() *InventoryComponent {
	component := InventoryComponent{
		items: map[int]types.Item{},
	}
	return &component
}

func (i *InventoryComponent) Give(item types.Item) error {
	if _, ok := i.items[item.ID()]; ok {
		return fmt.Errorf("Item %d already owned by entity", item.ID())
	}
	i.items[item.ID()] = item
	return nil
}

func (i *InventoryComponent) Take(item types.Item) error {
	if item, ok := i.items[item.ID()]; ok {
		delete(i.items, item.ID())
		return nil
	}

	return fmt.Errorf("item %d not found", item.ID())
}
