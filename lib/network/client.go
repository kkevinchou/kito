package network

import (
	"encoding/json"
	"fmt"
	"net"
)

type commandFrameFunc func() int

func defaultCommandFrameFunc() int {
	return -69
}

type Client struct {
	id           int
	connection   net.Conn
	messageQueue chan *Message

	commandFrameFunc commandFrameFunc
}

func baseClient() *Client {
	return &Client{
		id:               -1,
		messageQueue:     make(chan *Message, messageQueueBufferSize),
		commandFrameFunc: defaultCommandFrameFunc,
	}
}

func NewClient(id int, connection net.Conn) *Client {
	client := baseClient()
	client.id = id
	client.connection = connection
	go queueIncomingMessages(client.connection, client.messageQueue)
	return client
}

func (c *Client) SetCommandFrameFunction(f commandFrameFunc) {
	c.commandFrameFunc = f
}

func Connect(host, port, connectionType string) (*Client, error) {
	conn, err := net.Dial(connectionType, fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		fmt.Println(connectionType, fmt.Sprintf("%s:%s", host, port))
		fmt.Println("HI1")
		return nil, err
	}
	client := baseClient()
	client.connection = conn

	acceptMessage, err := readAcceptMessage(conn)
	if err != nil {
		return nil, err
	}
	client.id = acceptMessage.ID

	go queueIncomingMessages(conn, client.messageQueue)

	return client, nil
}

func (c *Client) PullIncomingMessages() []*Message {
	var messages []*Message
	for i := 0; i < len(c.messageQueue); i++ {
		select {
		case message := <-c.messageQueue:
			messages = append(messages, message)
		default:
			return messages
		}
	}
	return messages
}

func (c *Client) sendMessage(message *Message) error {
	encoder := json.NewEncoder(c.connection)
	err := encoder.Encode(message)
	if err != nil {
		return err
	}
	return nil
}

// SendMessage sends the message through the client
func (c *Client) SendMessage(messageType MessageType, subMessage interface{}) error {
	var bodyBytes []byte
	var err error
	if subMessage != nil {
		bodyBytes, err = json.Marshal(subMessage)
		if err != nil {
			return err
		}
	}

	msg := &Message{
		SenderID:     c.ID(),
		CommandFrame: c.commandFrameFunc(),
		MessageType:  messageType,
		Body:         bodyBytes,
	}

	return c.sendMessage(msg)
}

func (c *Client) SyncReceiveMessage() *Message {
	return <-c.messageQueue
}

func (c *Client) ID() int {
	return c.id
}

func readAcceptMessage(conn net.Conn) (*AcceptMessage, error) {
	message := Message{}
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&message)
	if err != nil {
		return nil, err
	}

	acceptMessage := AcceptMessage{}
	err = json.Unmarshal(message.Body, &acceptMessage)
	if err != nil {
		return nil, err
	}

	return &acceptMessage, nil
}
