package interfaces

import "github.com/kkevinchou/ant/lib/math/vector"

type ItemI interface {
	OwnedBy() int
	Owned() bool
	Id() int
	Position() vector.Vector
}

type InventoryI interface {
	Give(ItemI)
	Take(int) ItemI
}
