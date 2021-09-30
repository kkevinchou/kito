package singleton

import (
	"github.com/kkevinchou/kito/kito/inputbuffer"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/network"
)

// I hate this. Find a better place to put your data, yo
type Singleton struct {
	// client fields
	PlayerID int
	CameraID int

	// server fields
	InputBuffer                *inputbuffer.InputBuffer
	IncomingPlayerMessage      *network.Message
	IncomingPlayerInputMessage *network.InputMessage

	// Common
	PlayerInput  map[int]input.Input
	CommandFrame int
}

func NewSingleton() *Singleton {
	return &Singleton{
		PlayerInput: map[int]input.Input{},
		InputBuffer: inputbuffer.NewInputBuffer(settings.MaxInputBufferCommandFrames),
	}
}
