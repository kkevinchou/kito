package network

import (
	"encoding/json"
	"fmt"
	"net"
)

const ()

type Client struct {
	connection   net.Conn
	messageQueue chan *Message
}

func NewClient() *Client {
	return &Client{
		messageQueue: make(chan *Message, messageQueueBufferSize),
	}
}

func NewClientFromConnection(connection net.Conn) *Client {
	client := NewClient()
	client.connection = connection
	go queueIncomingMessages(client.connection, client.messageQueue)
	return client
}

func (c *Client) Connect(host, port, connectionType string) (*AcceptMessage, error) {
	conn, err := net.Dial(connectionType, fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, err
	}
	c.connection = conn

	acceptMessage, err := readAcceptMessage(conn)
	if err != nil {
		return nil, err
	}

	go queueIncomingMessages(conn, c.messageQueue)

	return acceptMessage, nil
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

func (c *Client) SendMessage(message *Message) error {
	encoder := json.NewEncoder(c.connection)
	err := encoder.Encode(message)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendWrappedMessage(senderID int, messageType MessageType, subMessage interface{}) error {
	bodyBytes, err := json.Marshal(subMessage)
	if err != nil {
		return err
	}

	msg := &Message{
		SenderID:    senderID,
		MessageType: messageType,
		Body:        bodyBytes,
	}

	return c.SendMessage(msg)
}

func (c *Client) SyncReceiveMessage() *Message {
	return <-c.messageQueue
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
