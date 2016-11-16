package components

import "github.com/kkevinchou/ant/interfaces"

type InventoryComponent struct {
	items map[int]interfaces.ItemI
}

func NewInventoryComponent() *InventoryComponent {
	component := InventoryComponent{}
	return &component
}

func (i *InventoryComponent) Give(item interfaces.ItemI) {
	i.items[item.Id()] = item
}

func (i *InventoryComponent) Take(id int) (interfaces.ItemI, bool) {
	if item, ok := i.items[id]; ok {
		delete(i.items, id)
		return item, true
	}

	return nil, false
}
