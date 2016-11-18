package item

import (
	"errors"
	"fmt"

	"github.com/kkevinchou/ant/interfaces"
)

type Manager struct {
	items      map[int]interfaces.Item
	ownedItems map[int]interfaces.Item
}

func (i *Manager) Register(item interfaces.Item) {
	i.items[item.Id()] = item
}

func (i *Manager) Random() (interfaces.Item, error) {
	for _, val := range i.items {
		return val, nil
	}
	return nil, errors.New("Could not get random item")
}

func (i *Manager) PickUp(item interfaces.Item) error {
	if item, ok := i.items[item.Id()]; ok {
		delete(i.items, item.Id())
		i.ownedItems[item.Id()] = item
		return nil
	}

	return fmt.Errorf("Could not pick up item with id %d", item.Id())
}

func (i *Manager) Drop(item interfaces.Item) error {
	if item, ok := i.ownedItems[item.Id()]; ok {
		delete(i.ownedItems, item.Id())
		i.items[item.Id()] = item
		return nil
	}
	return fmt.Errorf("Could not drop item with id %d", item.Id())
}

func NewManager() *Manager {
	return &Manager{
		items:      map[int]interfaces.Item{},
		ownedItems: map[int]interfaces.Item{},
	}
}
