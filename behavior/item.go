package behavior

import (
	"time"

	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/interfaces"
	"github.com/kkevinchou/ant/logger"
)

type PickupItem struct {
	Entity interfaces.ItemReceiver
}

func (p *PickupItem) Tick(input interface{}, state AIState, delta time.Duration) (interface{}, Status) {
	logger.Debug("PickupItem - ENTER")
	var item interfaces.Item
	var ok bool

	if item, ok = input.(interfaces.Item); !ok {
		logger.Debug("PickupItem - FAIL")
		return nil, FAILURE
	}

	itemManager := directory.GetDirectory().ItemManager()
	err := itemManager.PickUp(p.Entity, item)
	if err != nil {
		logger.Debug("PickupItem - FAIL")
		return nil, FAILURE
	}

	p.Entity.Give(item)
	logger.Debug("PickupItem - SUCCESS")
	return nil, SUCCESS
}

func (p *PickupItem) Reset() {}

type DropItem struct {
	Entity interfaces.ItemGiver
}

func (d *DropItem) Tick(input interface{}, state AIState, delta time.Duration) (interface{}, Status) {
	logger.Debug("DropItem - ENTER")

	var item interfaces.Item
	var ok bool

	if item, ok = input.(interfaces.Item); !ok {
		logger.Debug("DropItem - FAIL")
		return nil, FAILURE
	}

	itemManager := directory.GetDirectory().ItemManager()
	err := itemManager.Drop(d.Entity, item)
	if err != nil {
		logger.Debug("DropItem - FAIL")
		return nil, FAILURE
	}

	logger.Debug("DropItem - SUCCESS")
	return nil, SUCCESS
}

func (d *DropItem) Reset() {}

type RandomItem struct{}

func (r *RandomItem) Tick(input interface{}, state AIState, delta time.Duration) (interface{}, Status) {
	logger.Debug("RandomItem - ENTER")
	itemManager := directory.GetDirectory().ItemManager()
	item, err := itemManager.Random()
	if err != nil {
		logger.Debug("RandomItem - FAIL")
		return nil, FAILURE
	}
	return item, SUCCESS
}

func (r *RandomItem) Reset() {}
