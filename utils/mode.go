package utils

import "github.com/kkevinchou/kito/settings"

func IsClient() bool {
	return settings.CurrentGameMode == settings.GameModeClient
}

func IsServer() bool {
	return settings.CurrentGameMode == settings.GameModeServer
}
