package network

import (
	"encoding/json"
	"fmt"
	"net"
)

const ()

type Client struct {
	connection net.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(host, port, connectionType string) (*AcceptMessage, error) {
	conn, err := net.Dial(connectionType, fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, err
	}
	c.connection = conn

	message := Message{}
	decoder := json.NewDecoder(c.connection)
	err = decoder.Decode(&message)
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

func (c *Client) SendMessage() error {
	encoder := json.NewEncoder(c.connection)
	connectMessage := Message{
		SenderID:    0,
		MessageType: MessageTypeConnect,
	}

	err := encoder.Encode(connectMessage)
	if err != nil {
		return err
	}
	return nil
}
