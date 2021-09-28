package singleton

import (
	"github.com/kkevinchou/kito/kito/inputbuffer"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/lib/input"
)

// I hate this. Find a better place to put your data, yo
type Singleton struct {
	// client fields
	PlayerID int
	CameraID int

	// server fields
	*inputbuffer.InputBuffer

	// Common
	PlayerInput  map[int]input.Input
	CommandFrame int
}

func NewSingleton() *Singleton {
	return &Singleton{
		PlayerInput: map[int]input.Input{},
		InputBuffer: inputbuffer.NewInputBuffer(settings.InputBufferSize),
	}
}
