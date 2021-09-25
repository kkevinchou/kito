package behavior

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/behavior"
)

type Position struct {
	// TODO: write a test for this
	filler bool // empty structs share the same pointer address, this field prevents the node cache from accidentally caching
}

type Positionable interface {
	Position() mgl64.Vec3
}

func (p *Position) Tick(input interface{}, state behavior.AIState, delta time.Duration) (interface{}, behavior.Status) {
	if positionable, ok := input.(Positionable); ok {
		return positionable.Position(), behavior.SUCCESS
	}
	return nil, behavior.FAILURE
}

func (v *Position) Reset() {}
