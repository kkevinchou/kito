package network

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/input"
)

type MessageType int

const (
	MessageTypeConnect MessageType = iota
	MessageTypeAcceptConnection
	MessageTypeInput
	MessageTypeReplication
	MessageTypeCreatePlayer
	MessageTypeAckCreatePlayer
	MessageTypeGameStateUpdate
)

type Message struct {
	SenderID     int
	MessageType  MessageType
	CommandFrame int
	Timestamp    time.Time

	Body []byte
}

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
