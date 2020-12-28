package types

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
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
	Position() mgl64.Vec3
	SetPosition(position mgl64.Vec3)
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
	SetTarget(mgl64.Vec3)
	Velocity() mgl64.Vec3
	Heading() mgl64.Vec3
}

// Controllable is an entity that can be controlled
type Controllable interface {
	Controlled() bool
}

// Camera is a controllable entity who's perspective can be rendered from the render system
type Camera interface {
	Controllable

	Forward() mgl64.Vec3
	Right() mgl64.Vec3
	SetVelocity(vector mgl64.Vec3)
	MaxSpeed() float64

	Update(delta time.Duration)
	UpdateView(vector mgl64.Vec2)
	Position() mgl64.Vec3
	View() mgl64.Vec2
	ApplyImpulse(name string, impulse *Impulse)
}

type Singleton interface {
	GetKeyboardInputSet() *KeyboardInput
	SetKeyboardInputSet(input *KeyboardInput)
	GetMouseInput() *MouseInput
	SetMouseInput(input *MouseInput)
}
