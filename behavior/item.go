package behavior

import (
	"time"

	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/interfaces"
)

type PickupItem struct {
	Entity interfaces.ItemReceiver
}

func (p *PickupItem) Tick(input interface{}, state AIState, delta time.Duration) (interface{}, Status) {
	var item interfaces.Item
	var ok bool

	if item, ok = input.(interfaces.Item); !ok {
		return nil, FAILURE
	}

	itemManager := directory.GetDirectory().ItemManager()
	err := itemManager.PickUp(p.Entity, item)
	if err != nil {
		return nil, FAILURE
	}

	p.Entity.Give(item)
	return nil, SUCCESS
}

func (p *PickupItem) Reset() {}

type DropItem struct {
	Entity interfaces.ItemGiver
}

func (d *DropItem) Tick(input interface{}, state AIState, delta time.Duration) (interface{}, Status) {
	var item interfaces.Item
	var ok bool

	if item, ok = input.(interfaces.Item); !ok {
		return nil, FAILURE
	}

	itemManager := directory.GetDirectory().ItemManager()
	err := itemManager.Drop(d.Entity, item)
	if err != nil {
		return nil, FAILURE
	}

	return nil, SUCCESS
}

func (d *DropItem) Reset() {}

type RandomItem struct{}

func (r *RandomItem) Tick(input interface{}, state AIState, delta time.Duration) (interface{}, Status) {
	itemManager := directory.GetDirectory().ItemManager()
	item, err := itemManager.Random()
	if err != nil {
		return nil, FAILURE
	}
	return item, SUCCESS
}

func (r *RandomItem) Reset() {}
