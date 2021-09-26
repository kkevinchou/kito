package singleton

import (
	"github.com/kkevinchou/kito/lib/input"
)

type Singleton struct {
	// client fields
	PlayerID int
	CameraID int

	// server fields
	PlayerCommandFrames map[int]int // most recent command frame based on last input message

	// Common
	PlayerInput  map[int]input.Input
	CommandFrame int
}

func NewSingleton() *Singleton {
	return &Singleton{
		PlayerInput:         map[int]input.Input{},
		PlayerCommandFrames: map[int]int{},
	}
}
