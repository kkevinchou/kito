package inventory

type InventoryI interface {
	AddItem(interface{})
	DropItem(interface{})
	HasItem(interface{})
}

type ItemI interface {
}

type InventoryComponent struct {
}

func NewInventoryComponent() *InventoryComponent {
	component := InventoryComponent{}
	return &component
}

func (i *InventoryComponent) Give(item ItemI) {

}
