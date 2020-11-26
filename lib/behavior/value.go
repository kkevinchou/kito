package behavior

import "time"

type Value struct {
	Value interface{}
}

func (v *Value) Tick(input interface{}, state AIState, delta time.Duration) (interface{}, Status) {
	return v.Value, SUCCESS
}

func (v *Value) Reset() {}
