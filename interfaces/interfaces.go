package interfaces

import "github.com/kkevinchou/ant/lib/math/vector"

type Item interface {
	SetOwner(owner ItemReceiver)
	OwnedBy() ItemReceiver
	Owned() bool
	Position() vector.Vector
	Id() int
}

type ItemReceiver interface {
	Give(item Item) error
}

type ItemGiver interface {
	Take(item Item) error
}

type ItemGiverReceiver interface {
	ItemGiver
	ItemReceiver
}
