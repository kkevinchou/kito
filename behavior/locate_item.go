package behavior

import (
	"fmt"
	"time"

	"github.com/kkevinchou/ant/components"
	"github.com/kkevinchou/ant/systems"
)

type LocateItem struct {
	Entity components.InventoryI
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
