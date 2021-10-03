package types

import "github.com/kkevinchou/kito/lib/network"

type MovementType int

const (
	MovementTypeSteering    MovementType = iota
	MovementTypeDirectional MovementType = iota
)

type GameMode string

const (
	GameModeEditor  GameMode = "EDITOR"
	GameModePlaying GameMode = "PLAYING"
)

type NetworkClient interface {
	SendMessage(messageType int, messageBody interface{}) error
	PullIncomingMessages() []*network.Message
}
