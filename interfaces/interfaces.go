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

type Controllable interface {
	Forward() vector.Vector3
	Right() vector.Vector3
	SetVelocity(vector vector.Vector3)
	SetVelocityDirection(vector vector.Vector3)
	MaxSpeed() float64
}

type Worker interface {
	ItemGiverReceiver
	SetTarget(vector.Vector3)
	Velocity() vector.Vector3
	Heading() vector.Vector3
}

type Viewer interface {
	Controllable

	Position() vector.Vector3
	View() vector.Vector

	UpdateView(vector vector.Vector)

	Update(delta time.Duration)
}
