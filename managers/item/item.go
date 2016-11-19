package item

import (
	"errors"
	"fmt"

	"github.com/kkevinchou/ant/interfaces"
)

type Manager struct {
	items map[interfaces.Item]interface{}
}

func (i *Manager) Register(item interfaces.Item) {
	i.items[item] = nil
}

func (i *Manager) Random() (interfaces.Item, error) {
	for key, _ := range i.items {
		return key, nil
	}
	return nil, errors.New("Could not get random item")
}

func (i *Manager) PickUp(owner interfaces.ItemReceiver, item interfaces.Item) error {
	if _, ok := i.items[item]; ok {
		if err := owner.Give(item); err != nil {
			return err
		}
		item.SetOwner(owner)
		delete(i.items, item)
		return nil
	}

	return fmt.Errorf("Could not pick up item with id %d", item.ID())
}

func (i *Manager) Drop(owner interfaces.ItemGiver, item interfaces.Item) error {
	if _, ok := i.items[item]; !ok {
		if err := owner.Take(item); err != nil {
			return err
		}
		item.SetOwner(nil)
		item.SetPosition(owner.Position())
		i.items[item] = nil
		return nil
	}
	return fmt.Errorf("Could not drop item with id %d", item.ID())
}

func NewManager() *Manager {
	return &Manager{
		items: map[interfaces.Item]interface{}{},
	}
}
