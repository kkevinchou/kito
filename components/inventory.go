package components

import (
	"fmt"

	"github.com/kkevinchou/ant/interfaces"
)

type InventoryComponent struct {
	items map[int]interfaces.Item
}

func NewInventoryComponent() *InventoryComponent {
	component := InventoryComponent{
		items: map[int]interfaces.Item{},
	}
	return &component
}

func (i *InventoryComponent) Give(item interfaces.Item) error {
	if _, ok := i.items[item.ID()]; ok {
		return fmt.Errorf("Item %d already owned by entity", item.ID())
	}
	i.items[item.ID()] = item
	fmt.Println(len(i.items))
	return nil
}

func (i *InventoryComponent) Take(item interfaces.Item) error {
	if item, ok := i.items[item.ID()]; ok {
		delete(i.items, item.ID())
		fmt.Println(len(i.items))
		return nil
	}

	return fmt.Errorf("item %d not found", item.ID())
}
