package singleton

import (
	"github.com/kkevinchou/kito/lib/input"
)

type Singleton struct {
	// client fields
	PlayerID int
	CameraID int

	// Common
	PlayerInput  map[int]input.Input
	CommandFrame int
}

func NewSingleton() *Singleton {
	return &Singleton{
		PlayerInput: map[int]input.Input{},
	}
}
