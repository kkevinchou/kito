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

	// MaxInputBufferCommandFrames controls how many we buffer client inputs for. See the
	// InputBuffer struct definition for a more detailed description.

	// This buffer size should ideally be able to fully contain and fully sim a singular
	// player action. Otherwise, there's an edge case where a player starts an action and
	// the message takes more than a command frame of time to reach the server, causing the
	// player action to apply one or more frames late which causes a client misprediction.
	// If we were in the realm of the command buffer, we would have been able to place the
	// message in the buffer and "push" it forward and execute on the correct frame.

	// The max command frames for the buffer is currently static but is ideally dynamically
	// able to resize depending on the quality of the player's internet connection to the server.
	// The faster the connection, the smaller the buffer needs to be may be something we want to be
	// able to dyanmically adjust based on player latency. The larger their latency, the larger the
	// buffer

	// This is potentially overkill to avoiding absolutely no mispredictions on the client.
	// The drawback of an input buffer is we now add a delay before we process user inputs.
	MaxInputBufferCommandFrames int = 100 / int(MSPerCommandFrame)

	MaxStateBufferCommandFrames int = 100 / int(MSPerCommandFrame)

	// Aritificial latency for debugging purposes mostly.
	ArtificialClientLatency time.Duration = 0 * time.Millisecond
	ArtificialServerLatency time.Duration = 0 * time.Millisecond
)
