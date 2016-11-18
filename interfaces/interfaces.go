package interfaces

import "github.com/kkevinchou/ant/lib/math/vector"

type Item interface {
	OwnedBy() int
	Owned() bool
	Id() int
	Position() vector.Vector
}

type InventoryI interface {
	Give(Item)
	Take(int) Item
}

type ItemReceiver interface {
	Give(item Item)
}

type ItemGiver interface {
	Take(item Item)
}

type ItemGiverReceiver interface {
	ItemGiver
	ItemReceiver
}
