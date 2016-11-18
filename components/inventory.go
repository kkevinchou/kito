package components

import "github.com/kkevinchou/ant/interfaces"

type InventoryComponent struct {
	items map[int]interfaces.Item
}

func NewInventoryComponent() *InventoryComponent {
	component := InventoryComponent{}
	return &component
}

func (i *InventoryComponent) Give(item interfaces.Item) {
	i.items[item.Id()] = item
}

func (i *InventoryComponent) Take(id int) (interfaces.Item, bool) {
	if item, ok := i.items[id]; ok {
		delete(i.items, id)
		return item, true
	}

	return nil, false
}
