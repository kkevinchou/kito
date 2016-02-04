package behavior

import "time"

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

func (f *AddItem) Tick(state AIState, delta time.Duration) Status {
	return SUCCESS
}

func (f *DropItem) Tick(state AIState, delta time.Duration) Status {
	return SUCCESS
}

func (f *HaveItemCondition) Tick(state AIState, delta time.Duration) Status {
	return SUCCESS
}
