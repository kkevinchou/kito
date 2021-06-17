package network

type MessageType int

const (
	MessageTypeConnect MessageType = iota
	MessageTypeAcceptConnection
)

type Message struct {
	SenderID int `json:"sender_id`
	MessageType

	Body []byte `json:"body"`
}

type AcceptMessage struct {
	PlayerID int `json:"player_id"`
}
