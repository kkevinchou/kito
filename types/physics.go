package types

import (
	"time"

	"github.com/kkevinchou/kito/lib/math/vector"
)

type Impulse struct {
	Vector      vector.Vector3
	ElapsedTime time.Duration

	// the decay fraction per second for an impulse
	DecayRate float64
}
