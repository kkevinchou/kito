package singleton

import (
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/network"
)

type Singleton struct {
	// client fields
	PlayerID int
	Client   *network.Client

	Input input.Input
}

func NewSingleton() *Singleton {
	return &Singleton{}
}
