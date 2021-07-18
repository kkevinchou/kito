package settings

import "github.com/go-gl/mathgl/mgl64"

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
	Port           string = "8080"
	ConnectionType string = "tcp"

	GameModeUndefined GameMode = "UNDEFINED"
	GameModeClient    GameMode = "CLIENT"
	GameModeServer    GameMode = "SERVER"

	ServerID      = 69
	ServerIDStart = 70000
	EntityIDStart = 80000
)
