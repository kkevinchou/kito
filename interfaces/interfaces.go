package interfaces

import "github.com/kkevinchou/ant/lib/math/vector"

type Item interface {
	OwnedBy() int
	Owned() bool
	Id() int
	Position() vector.Vector
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
