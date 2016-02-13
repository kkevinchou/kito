package item

import (
	"errors"
	"fmt"

	"github.com/kkevinchou/ant/lib/math/vector"
)

type ItemI interface {
	OwnedBy() int
	Owned() bool
	Id() int
	Position() vector.Vector
}

type Manager struct {
	items      map[int]ItemI
	ownedItems map[int]ItemI
}

func (i *Manager) Register(item ItemI) {
	i.items[item.Id()] = item
}

func (i *Manager) Locate() (ItemI, error) {
	for _, val := range i.items {
		return val, nil
	}
	return nil, errors.New("Could not locate item")
}

func (i *Manager) PickUp(id int) (*ItemI, error) {
	if item, ok := i.items[id]; ok {
		delete(i.items, id)
		i.ownedItems[id] = item
		return &item, nil
	}
	return nil, fmt.Errorf("Could not pick up item with id %d", id)
}

func (i *Manager) Drop(id int) (*ItemI, error) {
	if item, ok := i.ownedItems[id]; ok {
		delete(i.ownedItems, id)
		i.items[id] = item
		return &item, nil
	}
	return nil, fmt.Errorf("Could not drop item with id %d", id)
}

func NewManager() *Manager {
	return &Manager{
		items:      map[int]ItemI{},
		ownedItems: map[int]ItemI{},
	}
}
