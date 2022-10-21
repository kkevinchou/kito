package singleton

import (
	"github.com/kkevinchou/kito/kito/inputbuffer"
	"github.com/kkevinchou/kito/kito/playercommand/protogen/playercommand"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/statebuffer"
	"github.com/kkevinchou/kito/lib/input"
)

// I hate this. Find a better place to put your data, yo
type Singleton struct {
	// client fields
	PlayerID    int
	CameraID    int
	StateBuffer *statebuffer.StateBuffer

	// server fields
	InputBuffer    *inputbuffer.InputBuffer
	PlayerCommands map[int]*playercommand.PlayerCommandList

	// Common
	PlayerInput  map[int]input.Input
	CommandFrame int
}

func NewSingleton() *Singleton {
	return &Singleton{
		PlayerInput:    map[int]input.Input{},
		PlayerCommands: map[int]*playercommand.PlayerCommandList{},
		StateBuffer:    statebuffer.NewStateBuffer(settings.MaxStateBufferCommandFrames),
		InputBuffer:    inputbuffer.NewInputBuffer(settings.MaxInputBufferCommandFrames),
	}
}
