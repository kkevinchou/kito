package behavior

import (
	"fmt"
	"time"

	"github.com/kkevinchou/ant/systems"
)

type ItemI interface {
	AddItem(interface{})
	DropItem(interface{})
	HasItem(interface{})
}

type AddItem struct {
	Entity ItemI
}

type DropItem struct {
	Entity ItemI
}

type HaveItemCondition struct {
	Entity ItemI
}

func (a *AddItem) Tick(state AiState, delta time.Duration) Status {
	return SUCCESS
}

func (d *DropItem) Tick(state AiState, delta time.Duration) Status {
	return SUCCESS
}

func (h *HaveItemCondition) Tick(state AiState, delta time.Duration) Status {
	return SUCCESS
}

type LocateItem struct {
	Entity ItemI
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
