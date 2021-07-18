package network

import (
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
	MessageTypeGameStateSnapshot
)

type Message struct {
	SenderID     int         `json:"sender_id"`
	MessageType  MessageType `json:"message_type"`
	CommandFrame int         `json:"command_frame"`

	Body []byte `json:"body"`
}

type AcceptMessage struct {
	ID int `json:"id"`
}

type CreatePlayerMessage struct {
}

type AckCreatePlayerMessage struct {
	ID          int        `json:"id"`
	CameraID    int        `json:"camera_id"`
	Position    mgl64.Vec3 `json:"transform"`
	Orientation mgl64.Quat `json:"orientation"`
}

type ReplicationMessage struct {
}

type EntitySnapshot struct {
	ID          int
	Position    mgl64.Vec3
	Orientation mgl64.Quat
}

type GameStateSnapshotMessage struct {
	Entities map[int]EntitySnapshot
}

type InputMessage struct {
	Input input.Input `json:"input"`
}
