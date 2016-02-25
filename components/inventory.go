package components

import (
	"github.com/kkevinchou/ant/interfaces"
	"github.com/kkevinchou/ant/lib/math/vector"
)

type NilItem struct {
}

func (n *NilItem) Id() int {
	return -1
}
func (n *NilItem) OwnedBy() int {
	return -1
}
func (n *NilItem) Owned() bool {
	return false
}
func (n *NilItem) Position() vector.Vector {
	return vector.Zero()
}

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

	return &NilItem{}, false
}
