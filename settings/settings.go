package settings

type GameMode string

var (
	CurrentGameMode GameMode = GameModeUndefined
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
)
