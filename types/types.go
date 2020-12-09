package types

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
