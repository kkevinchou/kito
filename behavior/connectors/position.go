package connectors

import (
	"time"

	"github.com/kkevinchou/ant/behavior"
	"github.com/kkevinchou/ant/lib/math/vector"
)

type Position struct{}

type Positionable interface {
	Position() vector.Vector
}

func (p *Position) Tick(input interface{}, state behavior.AIState, delta time.Duration) (interface{}, behavior.Status) {
	if positionable, ok := input.(Positionable); ok {
		return positionable.Position(), behavior.SUCCESS
	}
	return nil, behavior.FAILURE
}

func (v *Position) Reset() {}
