package settings

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
)

type GameMode string

var (
	CurrentGameMode GameMode = GameModeUndefined
)

var (
	CameraStartPosition = mgl64.Vec3{0, 10, 30}
	CameraStartView     = mgl64.Vec2{0, 0}
)

const (
	LoggingLevel          = 1
	Seed           int64  = 1234567
	Host           string = "localhost"
	RemoteHost     string = "localhost"
	Port           string = "8080"
	ConnectionType string = "tcp"

	GameModeUndefined GameMode = "UNDEFINED"
	GameModeClient    GameMode = "CLIENT"
	GameModeServer    GameMode = "SERVER"

	ServerID      int = 69
	ServerIDStart int = 70000
	EntityIDStart int = 80000

	// MSPerCommandFrame is the size of the simulation step for reading input,
	// physics, etc.
	MSPerCommandFrame float64 = 16
	FPS               float64 = 60

	// MaxTimeStepMS is the cap on how big a timestep on the game client can be.
	// The game world will probably fall apart (since we're losing time), but it
	// prevents sprials of death where the game falls further and further behind.
	MaxTimeStepMS float64 = 250 // in milliseconds

	// InputBufferSize controls how long we buffer client inputs for. See the
	// InputBuffer struct definition for a more detailed description.
	InputBufferSize int = 150 / int(MSPerCommandFrame)

	// Aritificial latency for debugging purposes mostly.
	ArtificialClientLatency time.Duration = 60 * time.Millisecond
)
