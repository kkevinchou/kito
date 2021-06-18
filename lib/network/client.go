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

func (c *Client) IncomingMessageQueue() chan *Message {
	return c.messageQueue
}

func (c *Client) SendMessage(message *Message) error {
	encoder := json.NewEncoder(c.connection)
	err := encoder.Encode(message)
	if err != nil {
		return err
	}
	return nil
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
