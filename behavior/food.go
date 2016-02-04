package behavior

import "time"

type FindFoodI interface {
}

type FindFood struct {
}

func (f *FindFood) Tick(state AIState, delta time.Duration) Status {
	return SUCCESS
}
