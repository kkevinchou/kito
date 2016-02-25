package behavior

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kkevinchou/ant/interfaces"
	"github.com/kkevinchou/ant/systems"
)

type PickupItem struct {
	Entity interfaces.InventoryI
}

func (p *PickupItem) Tick(state AiState, delta time.Duration) Status {
	itemId64, err := strconv.ParseInt(state.BlackBoard["output"], 10, 0)
	if err != nil {
		return FAILURE
	}
	itemId := int(itemId64)

	itemManager := systems.GetDirectory().ItemManager()
	item, err := itemManager.PickUp(itemId)
	if err != nil {
		return FAILURE
	}

	p.Entity.Give(item)
	return SUCCESS
}

type LocateItem struct {
}

// Locates a random item
func (l *LocateItem) Tick(state AiState, delta time.Duration) Status {
	itemManager := systems.GetDirectory().ItemManager()
	item, err := itemManager.Locate()
	if err != nil {
		return FAILURE
	}
	position := item.Position()
	state.BlackBoard["output"] = fmt.Sprintf("%f_%f", position.X, position.Y)
	return SUCCESS
}

func (l *LocateItem) Reset() {}
