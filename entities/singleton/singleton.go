package singleton

import (
	"github.com/kkevinchou/kito/lib/input"
)

type Singleton struct {
	// client fields
	PlayerID int

	// Common
	PlayerInput       map[int]input.Input
	CommandFrameCount int
}

func NewSingleton() *Singleton {
	return &Singleton{
		PlayerInput: map[int]input.Input{},
	}
}
