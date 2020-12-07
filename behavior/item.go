package behavior

import (
	"time"

	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/interfaces"
	"github.com/kkevinchou/kito/lib/behavior"
	"github.com/kkevinchou/kito/logger"
)

type PickupItem struct {
	Entity interfaces.ItemReceiver
}

func (p *PickupItem) Tick(input interface{}, state behavior.AIState, delta time.Duration) (interface{}, behavior.Status) {
	logger.Debug("PickupItem - ENTER")
	var item interfaces.Item
	var ok bool

	if item, ok = input.(interfaces.Item); !ok {
		logger.Debug("PickupItem - FAIL")
		return nil, behavior.FAILURE
	}

	itemManager := directory.GetDirectory().ItemManager()
	err := itemManager.PickUp(p.Entity, item)
	if err != nil {
		logger.Debug("PickupItem - FAIL")
		return nil, behavior.FAILURE
	}

	p.Entity.Give(item)
	logger.Debug("PickupItem - SUCCESS")
	return nil, behavior.SUCCESS
}

func (p *PickupItem) Reset() {}

type DropItem struct {
	Entity interfaces.ItemGiver
}

func (d *DropItem) Tick(input interface{}, state behavior.AIState, delta time.Duration) (interface{}, behavior.Status) {
	logger.Debug("DropItem - ENTER")

	var item interfaces.Item
	var ok bool

	if item, ok = input.(interfaces.Item); !ok {
		logger.Debug("DropItem - FAIL")
		return nil, behavior.FAILURE
	}

	itemManager := directory.GetDirectory().ItemManager()
	err := itemManager.Drop(d.Entity, item)
	if err != nil {
		logger.Debug("DropItem - FAIL")
		return nil, behavior.FAILURE
	}

	logger.Debug("DropItem - SUCCESS")
	return nil, behavior.SUCCESS
}

func (d *DropItem) Reset() {}

type RandomItem struct{}

func (r *RandomItem) Tick(input interface{}, state behavior.AIState, delta time.Duration) (interface{}, behavior.Status) {
	logger.Debug("RandomItem - ENTER")
	itemManager := directory.GetDirectory().ItemManager()
	item, err := itemManager.Random()
	if err != nil {
		logger.Debug("RandomItem - FAIL")
		return nil, behavior.FAILURE
	}
	return item, behavior.SUCCESS
}

func (r *RandomItem) Reset() {}
