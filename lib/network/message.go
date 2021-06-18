package network

type MessageType int

const (
	MessageTypeConnect MessageType = iota
	MessageTypeAcceptConnection
	MessageTypeInput
)

type Message struct {
	SenderID    int         `json:"sender_id"`
	MessageType MessageType `json:"message_type"`

	Body []byte `json:"body"`
}

type AcceptMessage struct {
	PlayerID int `json:"player_id"`
}
