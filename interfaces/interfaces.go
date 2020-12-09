package interfaces

import (
	"time"

	"github.com/kkevinchou/kito/kito/commands"
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/types"
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

// Controllable is an entity that can be controlled
type Controllable interface {
	Controlled() bool
}

// Viewer is a controllable entity who's perspective can be rendered from the render system
type Viewer interface {
	Controllable

	Forward() vector.Vector3
	Right() vector.Vector3
	SetVelocity(vector vector.Vector3)
	MaxSpeed() float64

	Update(delta time.Duration)
	UpdateView(vector vector.Vector)
	Position() vector.Vector3
	View() vector.Vector
	ApplyImpulse(name string, impulse *types.Impulse)
}

type Singleton interface {
	GetKeyboardInputSet() *commands.KeyboardInputSet
	SetKeyboardInputSet(input *commands.KeyboardInputSet)
	GetMouseInput() *types.MouseInput
	SetMouseInput(input *types.MouseInput)
}
