package interfaces

import (
	"time"

	"github.com/kkevinchou/kito/lib/math/vector"
)

type Item interface {
	Positionable
	IDable
	Ownable
}

type ItemReceiver interface {
	Give(item Item) error
}

type ItemGiver interface {
	Positionable
	Take(item Item) error
}

type ItemGiverReceiver interface {
	ItemGiver
	ItemReceiver
}

type Positionable interface {
	Position() vector.Vector3
	SetPosition(position vector.Vector3)
}

type IDable interface {
	ID() int
}

type Ownable interface {
	SetOwner(owner ItemReceiver)
	OwnedBy() ItemReceiver
	Owned() bool
}

type Worker interface {
	ItemGiverReceiver
	SetTarget(vector.Vector3)
	Velocity() vector.Vector3
	Heading() vector.Vector3
}

// Controllable is an entity that can be controlled. controlled via forward, backward, left, right, up, down
// and is set by SetVelocityDirection
type Controllable interface {
	// public
	SetVelocityDirection(vector vector.Vector3)

	// private
	Forward() vector.Vector3
	Right() vector.Vector3
	SetVelocity(vector vector.Vector3)
	MaxSpeed() float64
}

// Viewer is a controllable entity who's perspective can be rendered from the render system
type Viewer interface {
	Controllable

	// public
	Update(delta time.Duration)
	UpdateView(vector vector.Vector)
	Position() vector.Vector3
	View() vector.Vector
}
