package inventory

type InventoryI interface {
	AddItem(interface{})
	DropItem(interface{})
	HasItem(interface{})
}

type NilItem struct {
}

func (n *NilItem) Id() int {
	return -1
}

type ItemI interface {
	Id() int
}

type InventoryComponent struct {
	items map[int]ItemI
}

func NewInventoryComponent() *InventoryComponent {
	component := InventoryComponent{}
	return &component
}

func (i *InventoryComponent) Give(item ItemI) {
	i.items[item.Id()] = item
}

func (i *InventoryComponent) Take(id int) (ItemI, bool) {
	if item, ok := i.items[id]; ok {
		delete(i.items, id)
		return item, true
	}

	return &NilItem{}, false
}
