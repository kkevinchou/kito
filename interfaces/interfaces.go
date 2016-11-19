package interfaces

import "github.com/kkevinchou/ant/lib/math/vector"

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
	Position() vector.Vector
	SetPosition(position vector.Vector)
}

type IDable interface {
	ID() int
}

type Ownable interface {
	SetOwner(owner ItemReceiver)
	OwnedBy() ItemReceiver
	Owned() bool
}
