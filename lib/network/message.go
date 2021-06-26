package network

import "github.com/go-gl/mathgl/mgl64"

type MessageType int

const (
	MessageTypeConnect MessageType = iota
	MessageTypeAcceptConnection
	MessageTypeInput
	MessageTypeReplication
	MessageTypeCreatePlayer
	MessageTypeAckCreatePlayer
)

type Message struct {
	SenderID    int         `json:"sender_id"`
	MessageType MessageType `json:"message_type"`

	Body []byte `json:"body"`
}

type AcceptMessage struct {
	PlayerID int `json:"player_id"`
}

type CreatePlayerMessage struct {
}

type AckCreatePlayerMessage struct {
	Position    mgl64.Vec3 `json:"transform"`
	Orientation mgl64.Quat `json:"orientation"`
}

type ReplicationMessage struct {
}
