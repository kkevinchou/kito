package components

import (
	"github.com/kkevinchou/kito/kito/components/protogen/inventory"
	"google.golang.org/protobuf/proto"
)

type InventoryComponent struct {
	Data *inventory.Inventory
}

func NewInventoryComponent() *InventoryComponent {
	return &InventoryComponent{&inventory.Inventory{Items: []int64{}}}
}

func (c *InventoryComponent) AddToComponentContainer(container *ComponentContainer) {
	container.InventoryComponent = c
}

func (c *InventoryComponent) ComponentFlag() int {
	return ComponentFlagInventory
}

func (c *InventoryComponent) Synchronized() bool {
	return true
}

func (c *InventoryComponent) Load(bytes []byte) {
	h := &inventory.Inventory{}
	err := proto.Unmarshal(bytes, h)
	if err != nil {
		panic(err)
	}
	c.Data = h
}

func (c *InventoryComponent) Serialize() []byte {
	bytes, err := proto.Marshal(c.Data)
	if err != nil {
		panic(err)
	}
	return bytes
}
