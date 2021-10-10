package singleton

import (
	"github.com/kkevinchou/kito/kito/inputbuffer"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/statebuffer"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/metrics"
)

// I hate this. Find a better place to put your data, yo
type Singleton struct {
	// client fields
	PlayerID    int
	CameraID    int
	StateBuffer *statebuffer.StateBuffer

	// server fields
	InputBuffer *inputbuffer.InputBuffer

	// Common
	PlayerInput     map[int]input.Input
	CommandFrame    int
	MetricsRegistry *metrics.MetricsRegistry
}

func NewSingleton() *Singleton {
	return &Singleton{
		PlayerInput:     map[int]input.Input{},
		StateBuffer:     statebuffer.NewStateBuffer(settings.MaxStateBufferCommandFrames),
		InputBuffer:     inputbuffer.NewInputBuffer(settings.MaxInputBufferCommandFrames),
		MetricsRegistry: metrics.New(),
	}
}
