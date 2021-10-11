package knetwork

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/input"
)

const (
	MessageTypeAcceptConnection int = iota
	MessageTypeAckCreatePlayer
	MessageTypeInput
	MessageTypeCreatePlayer
	MessageTypeGameStateUpdate
	MessageTypePing
	MessageTypeAckPing
)

type AcceptMessage struct {
	ID int
}

type CreatePlayerMessage struct {
}

type AckCreatePlayerMessage struct {
	ID          int
	CameraID    int
	Position    mgl64.Vec3
	Orientation mgl64.Quat
}

type EntitySnapshot struct {
	ID          int
	Type        int
	Position    mgl64.Vec3
	Orientation mgl64.Quat
	Velocity    mgl64.Vec3
	Impulses    map[string]types.Impulse
}

type Event struct {
	Type  int
	Bytes []byte
}

type EventI interface {
	TypeAsInt() int
	Serialize() ([]byte, error)
}

type GameStateUpdateMessage struct {
	LastInputCommandFrame       int
	LastInputGlobalCommandFrame int
	CurrentGlobalCommandFrame   int
	Entities                    map[int]EntitySnapshot
	Events                      []Event
	Events2                     []EventI
}

type InputMessage struct {
	CommandFrame int
	Input        input.Input
}

type PingMessage struct {
	SendTime time.Time
}

type AckPingMessage struct {
	PingSendTime time.Time
}
